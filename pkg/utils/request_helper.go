package utils

import (
	"encoding/json"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"mime/multipart"
	"regexp"
	"strconv"
)

// ParseRegistrationRequest parses raw request body into ElectionPairRegistrationRequest
// This is in a separate file to avoid circular dependencies
func ParseRegistrationRequest(form *multipart.Form, files map[string]UploadedFile) (*request.ElectionPairRegistrationRequest, error) {
	pairValues := form.Value["pair"]
	if len(pairValues) == 0 {
		return nil, &custerr.ErrChain{
			Message: "Missing registration request data in 'pair' field",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	pairJSON := pairValues[0]
	var regReq request.ElectionPairRegistrationRequest
	if err := json.Unmarshal([]byte(pairJSON), &regReq); err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid JSON in registration request",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		}
	}

	if err := regReq.ValidateRegistrationRequest(); err != nil {
		return nil, err
	}

	if pairPhoto, exists := files["pair_photo"]; exists {
		regReq.PairPhotoName = pairPhoto.OriginalFilename
		regReq.PairPhotoFile = pairPhoto.File
	}

	if presidentPhoto, exists := files["president_photo"]; exists {
		regReq.President.PhotoName = presidentPhoto.OriginalFilename
		regReq.President.PhotoFile = presidentPhoto.File
	}

	if vicePresidentPhoto, exists := files["vice_president_photo"]; exists {
		regReq.VicePresident.PhotoName = vicePresidentPhoto.OriginalFilename
		regReq.VicePresident.PhotoFile = vicePresidentPhoto.File
	}

	return &regReq, nil
}

// ParseDetailRequest parses raw request body into ElectionPairDetailRequest
// This is in a separate file to avoid circular dependencies
func ParseUpsertDetailRequest(form *multipart.Form, files map[string]UploadedFile, workProgramPhotos map[string]UploadedFile) (*request.ElectionPairDetailRequest, error) {
	programValues := form.Value["detail"]
	if len(programValues) == 0 {
		return nil, &custerr.ErrChain{
			Message: "Missing detail request data in 'program' field",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	programJSON := programValues[0]
	var detailReq request.ElectionPairDetailRequest
	if err := json.Unmarshal([]byte(programJSON), &detailReq); err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid JSON in detail request",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		}
	}

	if err := detailReq.ValidateDetailRequest(); err != nil {
		return nil, err
	}

	if programDocs, exists := files["program_docs"]; exists {
		detailReq.ProgramDocsName = programDocs.OriginalFilename
		detailReq.ProgramDocsFile = programDocs.File
	}

	pattern := regexp.MustCompile(`^work_program_photo_(\d+)$`)
	for fieldName, fileInfo := range workProgramPhotos {
		matches := pattern.FindStringSubmatch(fieldName)
		if len(matches) < 2 {
			continue
		}

		index, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		if index < 0 || index >= len(detailReq.WorkProgram) {
			continue
		}

		detailReq.WorkProgram[index].ProgramPhoto = fileInfo.OriginalFilename
		detailReq.WorkProgram[index].ProgramPhotoFile = fileInfo.File
	}

	return &detailReq, nil
}
