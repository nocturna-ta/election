package controller

import (
	"context"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
)

// Ping godoc
// @Summary 	Ping
// @Description	Ping to check health
// @Tags		Health
// @Accept		json
// @Produce		json
// @Success		200	{object}	jsonResponse{}
// @Router		/health	[get]
func (api *API) Ping(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.Ping")
	defer span.End()

	return rest.NewJSONResponse().SetData("ok"), nil
}
