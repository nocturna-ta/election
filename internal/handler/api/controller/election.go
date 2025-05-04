package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/infrastructures/cutresp"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/http/filehandler"
	"github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
	"regexp"
	"strconv"
	"strings"
)

// RegisterElectionPair godoc
// @Summary 	Election
// @Description	Register Election Pair (President and Vice President)
// @Tags		Election
// @Accept 		multipart/form-data
// @Param 		X-User-Id header string false "Authorized User"
// @Param 		X-Address-Id header string false "Authorized Address"
// @Param 		X-Role header string false "Authorized Role"
// @Param 		pair formData string true  "Registration Request (JSON String)"
// @Param 		pair_photo formData file true "Pair Photo (jpg, jpeg, png only)"
// @Param 		president_photo formData file true "President Photo (jpg, jpeg, png only)"
// @Param 		vice_president_photo formData file true "Vice President Photo (jpg, jpeg, png only)"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.ElectionPairResponse}
// @Router		/v1/election/pairs/register	[post]
func (api *API) RegisterElectionPair(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.RegisterElectionPair")
	defer span.End()

	form, err := req.RawRequest().MultipartForm()
	if err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Failed to parse multipart form",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		})
	}
	fileConfigs := []utils.FileUploadConfig{
		{
			FieldName:  "pair_photo",
			Required:   true,
			UploadFunc: filehandler.ImageUploadOptions,
			ErrorMsgs: map[error]string{
				filehandler.ErrInvalidFileFormat: "Invalid file format for pair photo. Only JPG, JPEG, and PNG files are allowed",
			},
		},
		{
			FieldName:  "president_photo",
			Required:   true,
			UploadFunc: filehandler.ImageUploadOptions,
			ErrorMsgs: map[error]string{
				filehandler.ErrInvalidFileFormat: "Invalid file format for president photo. Only JPG, JPEG, and PNG files are allowed",
			},
		},
		{
			FieldName:  "vice_president_photo",
			Required:   true,
			UploadFunc: filehandler.ImageUploadOptions,
			ErrorMsgs: map[error]string{
				filehandler.ErrInvalidFileFormat: "Invalid file format for vice president photo. Only JPG, JPEG, and PNG files are allowed",
			},
		},
	}

	uploadedFiles, err := utils.ProcessFileUploads(ctx, form, fileConfigs)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	defer utils.CloseFiles(uploadedFiles)

	registrationRequest, err := utils.ParseRegistrationRequest(form, uploadedFiles)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	res, err := api.electionUc.RegisterElectionPair(ctx, registrationRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetElectionPairByNo godoc
// @Summary 	Election
// @Description	Get Election Pair By Number
// @Tags		Election
// @Accept		json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param 		no path string true "Election Number"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.ElectionPairResponse}
// @Router		/v1/election/pairs/number/{no}	[get]
func (api *API) GetElectionPairByNo(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetElectionPairByNo")
	defer span.End()

	no := req.Params("no")

	res, err := api.electionUc.GetElectionPairByNo(ctx, no)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetAllElectionPairs godoc
// @Summary 	Election
// @Description	Get All Election Pairs
// @Tags		Election
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Accept		json
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.ElectionPairListResponse}
// @Router		/v1/election/pairs	[get]
func (api *API) GetAllElectionPairs(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetAllElectionPairs")
	defer span.End()

	res, err := api.electionUc.GetAllElectionPairs(ctx)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetElectionPairByID godoc
// @Summary 	Election
// @Description	Get Election Pair By ID
// @Tags		Election
// @Accept		json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param 		id path string true "Election Pair ID"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.ElectionPairResponse}
// @Router		/v1/election/pairs/{id}	[get]
func (api *API) GetElectionPairByID(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetElectionPairByID")
	defer span.End()

	id := req.Params("id")

	res, err := api.electionUc.GetElectionPairByID(ctx, id)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetElectionPairDetail godoc
// @Summary 	Election Detail
// @Description	Get Election Pair Detail
// @Tags		Election-Detail
// @Accept		json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param 		pairID path string true "Election Pair ID"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.ElectionPairDetailResponse}
// @Router		/v1/election/pairs/{pairID}/detail	[get]
func (api *API) GetElectionPairDetail(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetElectionPairDetail")
	defer span.End()

	pairID := req.Params("pairID")

	res, err := api.electionUc.GetElectionPairDetail(ctx, pairID)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// UpsertElectionPairDetail godoc
// @Summary     Create or update election pair details with support for multiple work program photos
// @Description Create or Update Election Pair Detail with multiple work programs. For work program photos, use naming convention 'work_program_photo_[index]' where index matches the position in the work_program array.
// @Tags        Election-Detail
// @Accept      multipart/form-data
// @Param       X-User-Id header string false "Authorized User"
// @Param       X-Address-Id header string false "Authorized Address"
// @Param       X-Role header string false "Authorized Role"
// @Param       detail formData string true "Detail Request (JSON String)"
// @Param       program_docs formData file true "Program Documents (pdf, docx only)"
// @Param       work_program_photo_* formData file false "Photos for work programs (jpg, jpeg, png only). Use pattern work_program_photo_0, work_program_photo_1, etc."
// @Produce     json
// @Success     200 {object} jsonResponse{data=response.ElectionPairDetailResponse}
// @Router      /v1/election/pairs/detail [post]
func (api *API) UpsertElectionPairDetail(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.UpsertElectionPairDetail")
	defer span.End()

	// Get multipart form
	form, err := req.RawRequest().MultipartForm()
	if err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Failed to parse multipart form",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		})
	}

	// Get the detail JSON from the form
	detailValues := form.Value["detail"]
	if len(detailValues) == 0 {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Missing detail data",
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	// Parse the detail JSON
	var detailReq request.ElectionPairDetailRequest
	if err := json.Unmarshal([]byte(detailValues[0]), &detailReq); err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Invalid JSON in detail request",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		})
	}

	// Validate the detail request
	if err := detailReq.ValidateDetailRequest(); err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	// Process program docs
	if fileHeaders, ok := form.File["program_docs"]; ok && len(fileHeaders) > 0 {
		file, err := fileHeaders[0].Open()
		if err != nil {
			return cutresp.CustomErrorResponse(&custerr.ErrChain{
				Message: "Failed to open program docs file",
				Code:    400,
				Type:    response.ErrBadRequest,
				Cause:   err,
			})
		}
		defer file.Close()

		detailReq.ProgramDocsName = fileHeaders[0].Filename
		detailReq.ProgramDocsFile = file
	}

	// Process work program photos
	// This uses a regex to match field names like work_program_photo_0, work_program_photo_1, etc.
	for fieldName, fileHeaders := range form.File {
		if matched, _ := regexp.MatchString(`^work_program_photo_\d+$`, fieldName); matched && len(fileHeaders) > 0 {
			// Extract the index from the field name
			indexStr := strings.TrimPrefix(fieldName, "work_program_photo_")
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				continue // Skip if index cannot be parsed
			}

			// Check if the index is valid for the work program array
			if index < 0 || index >= len(detailReq.WorkProgram) {
				continue // Skip if index is out of bounds
			}

			// Open the file
			file, err := fileHeaders[0].Open()
			if err != nil {
				return cutresp.CustomErrorResponse(&custerr.ErrChain{
					Message: fmt.Sprintf("Failed to open work program photo %s", fieldName),
					Code:    400,
					Type:    response.ErrBadRequest,
					Cause:   err,
				})
			}
			defer file.Close()

			// Associate the file with the work program
			detailReq.WorkProgram[index].ProgramPhotoName = fileHeaders[0].Filename
			detailReq.WorkProgram[index].ProgramPhotoFile = file
		}
	}

	// Call the use case
	res, err := api.electionUc.UpsertElectionPairDetail(ctx, &detailReq)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetElectionPairPhoto godoc
// @Summary     Election Pair Photo
// @Description Get Election Pair Photo
// @Tags        Election-Images
// @Param X-User-Id header string false "User"
// @Param X-Address-Id header string false "Address"
// @Param X-Role header string false "Role"
// @Param       id path string true "Election Pair ID"
// @Produce     image/jpeg
// @Produce     image/png
// @Success     200
// @Router      /v1/election/pairs/{id}/photo [get]
func (api *API) GetElectionPairPhoto(ctx context.Context, req *router.Request) (*rest.AttachmentResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetElectionPairPhoto")
	defer span.End()

	id, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid Election Pair ID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	file, contentType, err := api.electionUc.GetElectionPairPhoto(ctx, id)
	if err != nil {
		return nil, err
	}

	return rest.NewAttachmentResponse().SetFile(file).SetFileName(file.FileName).SetContentType(contentType), nil
}

// GetPresidentPhoto godoc
// @Summary    President Photo
// @Description Get President Photo
// @Tags        Election-Images
// @Param X-User-Id header string false "User"
// @Param X-Address-Id header string false "Address"
// @Param X-Role header string false "Role"
// @Param       id path string true "Election Pair ID"
// @Produce     image/jpeg
// @Produce     image/png
// @Success     200
// @Router      /v1/election/pairs/{id}/photo/president [get]
func (api *API) GetPresidentPhoto(ctx context.Context, req *router.Request) (*rest.AttachmentResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetPresidentPhoto")
	defer span.End()

	id, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid Election Pair ID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	file, contentType, err := api.electionUc.GetPresidentPhoto(ctx, id)
	if err != nil {
		return nil, err
	}

	return rest.NewAttachmentResponse().SetFile(file).SetFileName(file.FileName).SetContentType(contentType), nil
}

// GetVicePresidentPhoto godoc
// @Summary    Vice President Photo
// @Description Get Vice President Photo
// @Tags        Election-Images
// @Param X-User-Id header string false "User"
// @Param X-Address-Id header string false "Address"
// @Param X-Role header string false "Role"
// @Param       id path string true "Election Pair ID"
// @Produce     image/jpeg
// @Produce     image/png
// @Success     200
// @Router      /v1/election/pairs/{id}/photo/vice-president [get]
func (api *API) GetVicePresidentPhoto(ctx context.Context, req *router.Request) (*rest.AttachmentResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetVicePresidentPhoto")
	defer span.End()

	id, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid Election Pair ID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	file, contentType, err := api.electionUc.GetVicePresidentPhoto(ctx, id)
	if err != nil {
		return nil, err
	}

	return rest.NewAttachmentResponse().SetFile(file).SetFileName(file.FileName).SetContentType(contentType), nil
}

// ActivateElectionPair godoc
// @Summary 	Election
// @Description	Activate Election Pair
// @Tags		Election
// @Accept		json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param 		activateElection body request.ElectionPairActivationRequest true "Activate Election Pair Request"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.ElectionPairActivationResponse}
// @Router		/v1/election/pairs/{id}/activate	[put]
func (api *API) ActivateElectionPair(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.ActivateElectionPair")
	defer span.End()

	var activateRequest request.ElectionPairActivationRequest
	err := json.Unmarshal(req.RawBody(), &activateRequest)

	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	res, err := api.electionUc.ActivateElectionPair(ctx, &activateRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil

}
