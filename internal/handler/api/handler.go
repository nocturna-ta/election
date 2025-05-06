package api

import (
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/election/internal/handler/api/controller"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/router"
)

type Options struct {
	Cfg               config.MainConfig
	ElectionUc        usecases.ElectionUseCases
	PartyUc           usecases.PartyUseCases
	SupportingPartyUc usecases.SupportingPartyUseCases
}

type Handler struct {
	opts        *Options
	listenErrCh chan error
	myRouter    *router.FastRouter
}

func New(opts *Options) *Handler {
	handler := &Handler{
		opts: opts,
	}
	handler.myRouter = controller.New(&controller.Options{
		Prefix:            opts.Cfg.API.BasePath,
		Port:              opts.Cfg.Server.Port,
		ReadTimeout:       opts.Cfg.Server.ReadTimeout,
		WriteTimeout:      opts.Cfg.Server.WriteTimeout,
		RequestTimeout:    opts.Cfg.API.APITimeout,
		EnableSwagger:     opts.Cfg.API.EnableSwagger,
		CorsConfig:        utils.ConvertToRouterCorsConfig(&opts.Cfg.Cors),
		ElectionUc:        opts.ElectionUc,
		PartyUc:           opts.PartyUc,
		SupportingPartyUc: opts.SupportingPartyUc,
	}).RegisterRoute()
	return handler
}
func (h *Handler) Run() {
	log.Infof("API Listening on %d", h.opts.Cfg.Server.Port)
	h.listenErrCh <- h.myRouter.StartServe()
}

func (h *Handler) ListenError() <-chan error {
	return h.listenErrCh
}
