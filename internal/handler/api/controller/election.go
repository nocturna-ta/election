package controller

import (
	"context"
	"encoding/json"
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
// @Tags		Election
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
// @Summary 	Election Detail
// @Description	Create or Update Election Pair Detail
// @Tags		Election
// @Accept		json
// @Param 		detail body request.ElectionPairDetailRequest true "Detail Request"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.ElectionPairDetailResponse}
// @Router		/v1/election/pairs/detail	[post]
func (api *API) UpsertElectionPairDetail(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.UpsertElectionPairDetail")
	defer span.End()

	var detailReq request.ElectionPairDetailRequest
	err := json.Unmarshal(req.RawBody(), &detailReq)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	if err := detailReq.ValidateDetailRequest(); err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	res, err := api.electionUc.UpsertElectionPairDetail(ctx, &detailReq)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetElectionPairPhoto godoc
// @Summary     Election Pair Photo
// @Description Get Election Pair Photo
// @Tags        Election
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
// @Tags        Election
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
// @Tags        Election
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
