package controller

import (
	"context"
	"github.com/nocturna-ta/election/internal/infrastructures/cutresp"
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
// @Accept		json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address-Id header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param 		pair body request.ElectionPairRegistrationRequest true "Registration Request"
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
		},
		{
			FieldName:  "vice_president_photo",
			Required:   true,
			UploadFunc: filehandler.ImageUploadOptions,
		},
	}

	uploadedFiles, err := utils.ProcessFileUploads(ctx, form, fileConfigs)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	defer utils.CloseFiles(uploadedFiles)

	registraionRequest, err := utils.ParseRegistrationRequest(req.RawBody(), uploadedFiles)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	res, err := api.electionUc.RegisterElectionPair(ctx, registraionRequest)
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
