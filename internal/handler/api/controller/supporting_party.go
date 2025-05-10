package controller

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/infrastructures/cutresp"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
)

// AddSupportingParty godoc
// @Summary 	Add Supporting Party
// @Description	Add a political party as a supporting party for an election pair
// @Tags		Supporting-Party
// @Accept		json
// @Param 		X-User-Id header string false "Authorized User"
// @Param 		X-Address-Id header string false "Authorized Address"
// @Param 		X-Role header string false "Authorized Role"
// @Param 		supportingParty body request.AddSupportingPartyRequest true "Add Supporting Party Request"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.SupportingPartyResponse}
// @Router		/v1/election/pairs/supporting-party [post]
func (api *API) AddSupportingParty(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.AddSupportingParty")
	defer span.End()

	var addRequest request.AddSupportingPartyRequest
	err := json.Unmarshal(req.RawBody(), &addRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Invalid request format",
			Cause:   err,
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	if err := addRequest.Validate(); err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Invalid request data",
			Cause:   err,
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	res, err := api.supportingPartyUc.AddSupportingParty(ctx, &addRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetSupportingPartiesByPairID godoc
// @Summary 	Get Supporting Parties
// @Description	Get all supporting parties for an election pair
// @Tags		Supporting-Party
// @Accept		json
// @Param 		X-User-Id header string false "Authorized User"
// @Param 		X-Address-Id header string false "Authorized Address"
// @Param 		X-Role header string false "Authorized Role"
// @Param 		pairID path string true "Election Pair ID"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.SupportingPartiesResponse}
// @Router		/v1/election/pairs/{pairID}/supporting-parties [get]
func (api *API) GetSupportingPartiesByPairID(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetSupportingPartiesByPairID")
	defer span.End()

	pairID, err := uuid.Parse(req.Params("pairID"))
	if err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Invalid election pair ID format",
			Cause:   err,
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	res, err := api.supportingPartyUc.GetSupportingPartiesByPairID(ctx, pairID)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// RemoveSupportingParty godoc
// @Summary 	Remove Supporting Party
// @Description	Remove a political party from supporting an election pair
// @Tags		Supporting-Party
// @Accept		json
// @Param 		X-User-Id header string false "Authorized User"
// @Param 		X-Address-Id header string false "Authorized Address"
// @Param 		X-Role header string false "Authorized Role"
// @Param 		supportingParty body request.RemoveSupportingPartyRequest true "Remove Supporting Party Request"
// @Produce		json
// @Success		200	{object}	jsonResponse{}
// @Router		/v1/election/pairs/supporting-party [delete]
func (api *API) RemoveSupportingParty(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.RemoveSupportingParty")
	defer span.End()

	var removeRequest request.RemoveSupportingPartyRequest
	err := json.Unmarshal(req.RawBody(), &removeRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Invalid request format",
			Cause:   err,
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	if err := removeRequest.Validate(); err != nil {
		return cutresp.CustomErrorResponse(&custerr.ErrChain{
			Message: "Invalid request data",
			Cause:   err,
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	err = api.supportingPartyUc.RemoveSupportingParty(ctx, &removeRequest)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetMessage("Supporting party removed successfully"), nil
}
