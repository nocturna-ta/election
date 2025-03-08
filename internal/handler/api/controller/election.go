package controller

import (
	"context"
	"encoding/json"
	"github.com/nocturna-ta/election/internal/infrastructures/cutresp"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
)

// RegisterCandidate godoc
// @Summary 	Election
// @Description	Register Candidate
// @Tags		Election
// @Accept		json
// @Param 		candidate body request.CandidateRegistrationRequest true "Register Request"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.CandidateResponse}
// @Router		/v1/election/register	[post]
func (api *API) RegisterCandidate(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.RegisterCandidate")
	defer span.End()

	var regisReq request.CandidateRegistrationRequest
	err := json.Unmarshal(req.RawBody(), &regisReq)

	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	err = regisReq.ValidateRegistrationRequest()

	res, err := api.electionUc.RegisterCandidate(ctx, &regisReq)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetAllCandidate godoc
// @Summary 	Election
// @Description	Get All Candidate
// @Tags		Election
// @Accept		json
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.CandidateResponse}
// @Router		/v1/election/candidates	[get]
func (api *API) GetAllCandidate(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetAllCandidate")
	defer span.End()

	res, err := api.electionUc.GetAllCandidate(ctx)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetCandidateByNo godoc
// @Summary 	Election
// @Description	Get Candidate By No
// @Tags		Election
// @Accept		json
// @Param 		no path string true "Election No"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.CandidateResponse}
// @Router		/v1/election/candidate/{no}	[get]
func (api *API) GetCandidateByNo(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetCandidateByNo")
	defer span.End()

	no := req.Params("no")

	res, err := api.electionUc.GetCandidateByNo(ctx, no)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// ActivateCandidate godoc
// @Summary 	Election
// @Description	Activate Candidate
// @Tags		Election
// @Accept		json
// @Param 		candidate body request.CandidateActivationRequest true "Activate Request"
// @Produce		json
// @Success		200	{object}	jsonResponse{data=response.CandidateActivation}
// @Router		/v1/election/activate	[post]
func (api *API) ActivateCandidate(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.ActivateCandidate")
	defer span.End()

	var actReq request.CandidateActivationRequest
	err := json.Unmarshal(req.RawBody(), &actReq)

	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	err = actReq.ValidateActivationRequest()

	res, err := api.electionUc.ActivateCandidate(ctx, &actReq)
	if err != nil {
		return cutresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}
