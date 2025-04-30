package controller

import (
	"github.com/gofiber/swagger"
	_ "github.com/nocturna-ta/election/docs"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/golib/router"
	"time"
)

type API struct {
	prefix         string
	port           uint
	readTimeout    time.Duration
	writeTimeout   time.Duration
	requestTimeout time.Duration
	enableSwagger  bool
	corsConfig     *router.CorsConfig
	electionUc     usecases.ElectionUseCases
}

type Options struct {
	Prefix         string
	Port           uint
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestTimeout time.Duration
	EnableSwagger  bool
	CorsConfig     *router.CorsConfig
	ElectionUc     usecases.ElectionUseCases
}

func New(opts *Options) *API {
	return &API{
		prefix:         opts.Prefix,
		port:           opts.Port,
		readTimeout:    opts.ReadTimeout,
		writeTimeout:   opts.WriteTimeout,
		requestTimeout: opts.RequestTimeout,
		enableSwagger:  opts.EnableSwagger,
		corsConfig:     opts.CorsConfig,
		electionUc:     opts.ElectionUc,
	}
}

func (api *API) RegisterRoute() *router.FastRouter {
	corsConfig := &router.CorsConfig{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE",
		AllowHeaders:     "Content-Type, Authorization, x-user-id",
		AllowCredentials: false,
		MaxAge:           300,
	}
	myRouter := router.New(&router.Options{
		Prefix:         api.prefix,
		Port:           api.port,
		ReadTimeout:    api.readTimeout,
		WriteTimeout:   api.writeTimeout,
		RequestTimeout: api.requestTimeout,
		CorsConfig:     corsConfig,
	})

	if api.enableSwagger {
		myRouter.CustomHandler("GET", "/docs/*", swagger.HandlerDefault, router.MustAuthorized(false))
	}

	myRouter.GET("/health", api.Ping, router.MustAuthorized(false))

	myRouter.Group("/v1", func(v1 *router.FastRouter) {
		v1.Group("/election", func(election *router.FastRouter) {
			election.Group("/pairs", func(pairs *router.FastRouter) {
				pairs.GET("", api.GetAllElectionPairs, router.MustAuthorized(false))
				pairs.GET("/:id", api.GetElectionPairByID, router.MustAuthorized(false))
				pairs.GET("/number/:no", api.GetElectionPairByNo, router.MustAuthorized(false))
				pairs.POST("/register", api.RegisterElectionPair, router.MustAuthorized(false))
				//pairs.POST("/activate", api.ActivateElectionPair, router.MustAuthorized(false))
				//pairs.GET("/:id/full", api.GetElectionPairFull, router.MustAuthorized(false))
			})
		})
	})
	return myRouter
}
