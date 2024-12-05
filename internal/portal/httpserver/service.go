package httpserver

import (
	"context"
	"net/http"
	"vc/internal/portal/apiv1"
	"vc/pkg/httphelpers"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"

	"github.com/gin-gonic/gin"
)

// Service is the service object for httpserver
type Service struct {
	cfg         *model.Cfg
	log         *logger.Log
	server      *http.Server
	apiv1       Apiv1
	gin         *gin.Engine
	tracer      *trace.Tracer
	httpHelpers *httphelpers.Client
}

// New creates a new httpserver service
func New(ctx context.Context, cfg *model.Cfg, apiv1 *apiv1.Client, tracer *trace.Tracer, log *logger.Log) (*Service, error) {
	s := &Service{
		cfg:    cfg,
		log:    log.New("httpserver"),
		apiv1:  apiv1,
		gin:    gin.New(),
		tracer: tracer,
		server: &http.Server{},
	}

	var err error
	s.httpHelpers, err = httphelpers.New(ctx, s.tracer, s.cfg, s.log)
	if err != nil {
		return nil, err
	}

	// extra middlewares (must be declared before Server.Default)
	s.gin.Use(s.httpHelpers.Middleware.Gzip(ctx))

	rgRoot, err := s.httpHelpers.Server.Default(ctx, s.server, s.gin, s.cfg.Portal.APIServer.Addr)
	if err != nil {
		return nil, err
	}

	s.gin.Static("/static", "./static")
	s.gin.LoadHTMLFiles("./static/index.html")
	s.gin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	s.httpHelpers.Server.RegEndpoint(ctx, rgRoot, http.MethodGet, "health", s.endpointHealth)

	//TODO(mk): impl auth and replace below with //rgSecure := rgRoot.Group("secure", s.middlewareAuthRequired(ctx))
	rgSecure := rgRoot.Group("secure")

	rgAPIGWSecure := rgSecure.Group("apigw")
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIGWSecure, http.MethodPost, "document/search", s.endpointSearchDocuments)

	// Run http server
	go func() {
		err := s.httpHelpers.Server.ListenAndServe(ctx, s.server, s.cfg.Portal.APIServer)
		if err != nil {
			s.log.Trace("listen_error", "error", err)
		}
	}()

	s.log.Info("Started")

	return s, nil
}

// Close closing httpserver
func (s *Service) Close(ctx context.Context) error {
	s.log.Info("Stopped")
	return nil
}
