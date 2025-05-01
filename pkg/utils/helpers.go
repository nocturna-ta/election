package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/fileutils"
	"github.com/nocturna-ta/golib/http/filehandler"
	"github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/router"
	"io"
	"mime/multipart"
)

type FileUploadConfig struct {
	FieldName   string
	Required    bool
	UploadFunc  func() *filehandler.UploadOptions
	ErrorMsgs   map[error]string
	DefaultCode int
}

type UploadedFile struct {
	File             io.ReadCloser
	OriginalFilename string
	FilePath         string
	ContentType      string
	Size             int64
}

// Generic registration request interface to avoid direct import of request package
type RegistrationRequest interface {
	Validate() error
}

func ProcessFileUploads(ctx context.Context, form *multipart.Form, configs []FileUploadConfig) (map[string]UploadedFile, error) {
	result := make(map[string]UploadedFile)

	for _, config := range configs {
		uploadOptions := config.UploadFunc()
		uploadOptions.FieldName = config.FieldName

		uploadResult, err := filehandler.UploadFile(ctx, form, uploadOptions)
		if err != nil {
			if errors.Is(err, filehandler.ErrNoFile) && !config.Required {
				continue
			}
			return nil, MapFileUploadError(err, config)
		}

		file, err := fileutils.OpenFile(ctx, uploadResult.FilePath)
		if err != nil {
			return nil, &custerr.ErrChain{
				Message: fmt.Sprintf("Failed to open uploaded file: %s", config.FieldName),
				Code:    500,
				Type:    response.ErrInternalServerError,
				Cause:   err,
			}
		}

		result[config.FieldName] = UploadedFile{
			File:             file,
			OriginalFilename: uploadResult.OriginalFilename,
			FilePath:         uploadResult.FilePath,
			ContentType:      uploadResult.ContentType,
			Size:             uploadResult.Size,
		}
	}

	return result, nil
}

func MapFileUploadError(err error, config FileUploadConfig) *custerr.ErrChain {
	var errorMsg string
	var errorCode int = config.DefaultCode
	if errorCode == 0 {
		errorCode = 400
	}

	if config.ErrorMsgs != nil {
		if msg, exists := config.ErrorMsgs[err]; exists {
			errorMsg = msg
		}
	}

	if errorMsg == "" {
		switch err {
		case filehandler.ErrNoFile:
			errorMsg = fmt.Sprintf("No file provided for %s", config.FieldName)
		case filehandler.ErrFileTooLarge:
			errorMsg = fmt.Sprintf("File size exceeds maximum allowed size for %s", config.FieldName)
		case filehandler.ErrInvalidFileFormat:
			errorMsg = fmt.Sprintf("Invalid file format for %s", config.FieldName)
		default:
			errorMsg = fmt.Sprintf("Failed to process uploaded file: %s", config.FieldName)
			errorCode = 500
		}
	}

	return &custerr.ErrChain{
		Message: errorMsg,
		Code:    errorCode,
		Type:    response.ErrBadRequest,
		Cause:   err,
	}
}

// This function is now moved to pkg/utils/request_helper.go to break the dependency cycle
// Moving it to a separate file that will depend on request, but not be imported by request

func CloseFiles(files map[string]UploadedFile) {
	for _, fileInfo := range files {
		if fileInfo.File != nil {
			fileInfo.File.Close()
		}
	}
}

func CloseReadClosers(closers ...io.ReadCloser) {
	for _, closer := range closers {
		if closer != nil {
			closer.Close()
		}
	}
}

func StringToTx(signedTx string) (*types.Transaction, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(common.FromHex(signedTx)); err != nil {
		return nil, err
	}

	return tx, nil
}

func ConvertToRouterCorsConfig(configCors *config.CorsConfig) *router.CorsConfig {
	return &router.CorsConfig{
		AllowOrigins:     configCors.AllowOrigins,
		AllowMethods:     configCors.AllowMethods,
		AllowHeaders:     configCors.AllowHeaders,
		AllowCredentials: configCors.AllowCredentials,
		ExposeHeaders:    configCors.ExposeHeaders,
		MaxAge:           configCors.MaxAge,
	}
}
