package controller

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/infrastructures/cutresp"
	"github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/http/filehandler"
	"github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
)

// RegisterParty godoc
// @Summary Party Registration
// @Description Register a new political party with logo upload
// @Tags Party
// @Accept multipart/form-data
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param party formData string true "Party Registration Request (JSON String)"
// @Param logo formData file false "Party Logo (jpg, jpeg, png only)"
// @Produce json
// @Success 200 {object} jsonResponse{data=response.PartyResponse}
// @Router /v1/party/register [post]
func (api *API) RegisterParty(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.RegisterParty")
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
			FieldName:  "logo",
			Required:   false,
			UploadFunc: filehandler.ImageUploadOptions,
			ErrorMsgs: map[error]string{
				filehandler.ErrInvalidFileFormat: "Invalid file format for logo. Only JPG, JPEG, and PNG files are allowed",
			},
		},
	}

	uploadedFiles, err := utils.ProcessFileUploads(ctx, form, fileConfigs)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	defer utils.CloseFiles(uploadedFiles)

	partyRequest, err := utils.ParsePartyRequest(form, uploadedFiles, false)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	res, err := api.partyUc.RegisterParty(ctx, partyRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// UpdateParty godoc
// @Summary Party Update
// @Description Update a political party with optional logo upload
// @Tags Party
// @Accept multipart/form-data
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param party formData string true "Party Update Request (JSON String)"
// @Param logo formData file false "Party Logo (jpg, jpeg, png only)"
// @Produce json
// @Success 200 {object} jsonResponse{data=response.PartyResponse}
// @Router /v1/party/update [put]
func (api *API) UpdateParty(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.UpdateParty")
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
			FieldName:  "logo",
			Required:   false,
			UploadFunc: filehandler.ImageUploadOptions,
			ErrorMsgs: map[error]string{
				filehandler.ErrInvalidFileFormat: "Invalid file format for logo. Only JPG, JPEG, and PNG files are allowed",
			},
		},
	}

	uploadedFiles, err := utils.ProcessFileUploads(ctx, form, fileConfigs)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	defer utils.CloseFiles(uploadedFiles)

	updateRequest, err := utils.ParsePartyUpdateRequest(form, uploadedFiles)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	res, err := api.partyUc.UpdateParty(ctx, updateRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// DeleteParty godoc
// @Summary Party Deletion
// @Description Delete a political party by ID
// @Tags Party
// @Accept json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param id path string true "Party ID"
// @Produce json
// @Success 200 {object} jsonResponse{}
// @Router /v1/party/{id} [delete]
func (api *API) DeleteParty(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.DeleteParty")
	defer span.End()

	id := req.Params("id")
	if id == "" {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Party ID is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	err := api.partyUc.DeleteParty(ctx, id)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetMessage("Party deleted successfully"), nil
}

// GetPartyByID godoc
// @Summary Get Party
// @Description Get political party by ID
// @Tags Party
// @Accept json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param id path string true "Party ID"
// @Produce json
// @Success 200 {object} jsonResponse{data=response.PartyResponse}
// @Router /v1/party/{id} [get]
func (api *API) GetPartyByID(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetPartyByID")
	defer span.End()

	id := req.Params("id")
	if id == "" {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Party ID is required",
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	res, err := api.partyUc.GetPartyByID(ctx, id)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetAllParties godoc
// @Summary Get All Parties
// @Description Get all political parties
// @Tags Party
// @Accept json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Produce json
// @Success 200 {object} jsonResponse{data=[]response.PartyResponse}
// @Router /v1/party [get]
func (api *API) GetAllParties(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetAllParties")
	defer span.End()

	res, err := api.partyUc.GetAllParties(ctx)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetPartyPhoto godoc
// @Summary    Party Photo
// @Description Get Party Photo
// @Tags        Party-Images
// @Param X-User-Id header string false "User"
// @Param X-Address-Id header string false "Address"
// @Param X-Role header string false "Role"
// @Param       id path string true "Party Pair ID"
// @Produce     image/jpeg
// @Produce     image/png
// @Success     200
// @Router      /v1/party/{id}/photo [get]
func (api *API) GetPartyPhoto(ctx context.Context, req *router.Request) (*rest.AttachmentResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetPartyPhoto")
	defer span.End()

	id, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid Party Pair ID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	file, contentType, err := api.partyUc.GetPartyPhoto(ctx, id)
	if err != nil {
		return nil, err
	}

	return rest.NewAttachmentResponse().SetFile(file).SetFileName(file.FileName).SetContentType(contentType), nil
}
