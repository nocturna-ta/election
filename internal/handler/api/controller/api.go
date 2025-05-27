package controller

import (
	"github.com/gofiber/swagger"
	_ "github.com/nocturna-ta/election/docs"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/golib/router"
	"time"
)

type API struct {
	prefix            string
	port              uint
	readTimeout       time.Duration
	writeTimeout      time.Duration
	requestTimeout    time.Duration
	enableSwagger     bool
	corsConfig        *router.CorsConfig
	electionUc        usecases.ElectionUseCases
	partyUc           usecases.PartyUseCases
	supportingPartyUc usecases.SupportingPartyUseCases
}

type Options struct {
	Prefix            string
	Port              uint
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	RequestTimeout    time.Duration
	EnableSwagger     bool
	CorsConfig        *router.CorsConfig
	ElectionUc        usecases.ElectionUseCases
	PartyUc           usecases.PartyUseCases
	SupportingPartyUc usecases.SupportingPartyUseCases
}

func New(opts *Options) *API {
	return &API{
		prefix:            opts.Prefix,
		port:              opts.Port,
		readTimeout:       opts.ReadTimeout,
		writeTimeout:      opts.WriteTimeout,
		requestTimeout:    opts.RequestTimeout,
		enableSwagger:     opts.EnableSwagger,
		corsConfig:        opts.CorsConfig,
		electionUc:        opts.ElectionUc,
		partyUc:           opts.PartyUc,
		supportingPartyUc: opts.SupportingPartyUc,
	}
}

func (api *API) RegisterRoute() *router.FastRouter {

	myRouter := router.New(&router.Options{
		Prefix:         api.prefix,
		Port:           api.port,
		ReadTimeout:    api.readTimeout,
		WriteTimeout:   api.writeTimeout,
		RequestTimeout: api.requestTimeout,
		CorsConfig:     api.corsConfig,
	})

	if api.enableSwagger {
		myRouter.CustomHandler("GET", "/docs/*", swagger.HandlerDefault, router.MustAuthorized(false))
	}

	myRouter.GET("/health", api.Ping, router.MustAuthorized(false))

	myRouter.Group("/v1", func(v1 *router.FastRouter) {
		v1.Group("/election", func(election *router.FastRouter) {
			election.Group("/pairs", func(pairs *router.FastRouter) {
				pairs.Group("/:id", func(id *router.FastRouter) {
					id.GET("", api.GetElectionPairByID, router.MustAuthorized(false))
					id.ATTACHMENT("/photo", api.GetElectionPairPhoto, router.MustAuthorized(false))
					id.ATTACHMENT("/photo/president", api.GetPresidentPhoto, router.MustAuthorized(false))
					id.ATTACHMENT("/photo/vice-president", api.GetVicePresidentPhoto, router.MustAuthorized(false))
					id.PUT("/activate", api.ActivateElectionPair, router.MustAuthorized(false))
					id.GET("/full", api.GetElectionPairFull, router.MustAuthorized(false))
				})
				pairs.GET("", api.GetAllElectionPairs, router.MustAuthorized(false))
				pairs.GET("/number/:no", api.GetElectionPairByNo, router.MustAuthorized(false))
				pairs.POST("/register", api.RegisterElectionPair, router.MustAuthorized(false))
				pairs.POST("/detail", api.UpsertElectionPairDetail, router.MustAuthorized(false))
				pairs.GET("/:pairID/detail", api.GetElectionPairDetail, router.MustAuthorized(false))

				pairs.POST("/supporting-party", api.AddSupportingParty, router.MustAuthorized(false))
				pairs.DELETE("/supporting-party", api.RemoveSupportingParty, router.MustAuthorized(false))
				pairs.GET("/:pairID/supporting-parties", api.GetSupportingPartiesByPairID, router.MustAuthorized(false))
			})
		})

		v1.Group("/party", func(party *router.FastRouter) {
			party.POST("/register", api.RegisterParty, router.MustAuthorized(false))
			party.PUT("/update", api.UpdateParty, router.MustAuthorized(false))
			party.GET("/:id", api.GetPartyByID, router.MustAuthorized(false))
			party.DELETE("/:id", api.DeleteParty, router.MustAuthorized(false))
			party.ATTACHMENT("/:id/photo", api.GetPartyPhoto, router.MustAuthorized(false))
			party.GET("", api.GetAllParties, router.MustAuthorized(false))
		})
	})
	return myRouter
}
