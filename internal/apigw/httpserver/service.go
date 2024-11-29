package httpserver

import (
	"context"
	"net/http"
	"vc/internal/apigw/apiv1"
	"vc/pkg/httphelpers"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"

	// Swagger
	_ "vc/docs/apigw"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Service is the service object for httpserver
type Service struct {
	cfg            *model.Cfg
	log            *logger.Log
	server         *http.Server
	apiv1          Apiv1
	gin            *gin.Engine
	tracer         *trace.Tracer
	eventPublisher apiv1.EventPublisher
	httpHelpers    *httphelpers.Client
}

// New creates a new httpserver service
func New(ctx context.Context, cfg *model.Cfg, apiv1 *apiv1.Client, tracer *trace.Tracer, eventPublisher apiv1.EventPublisher, log *logger.Log) (*Service, error) {
	s := &Service{
		cfg:            cfg,
		log:            log.New("httpserver"),
		apiv1:          apiv1,
		gin:            gin.New(),
		tracer:         tracer,
		server:         &http.Server{},
		eventPublisher: eventPublisher,
	}

	var err error
	s.httpHelpers, err = httphelpers.New(ctx, s.tracer, s.cfg, s.log)
	if err != nil {
		return nil, err
	}

	rgRoot, err := s.httpHelpers.Server.Default(ctx, s.server, s.gin, s.cfg.APIGW.APIServer.Addr)
	if err != nil {
		return nil, err
	}

	s.httpHelpers.Server.RegEndpoint(ctx, rgRoot, http.MethodGet, "health", s.endpointHealth)

	rgDocs := rgRoot.Group("/swagger")
	rgDocs.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rgAPIv1 := rgRoot.Group("api/v1")

	if s.cfg.APIGW.APIServer.BasicAuth.Enabled {
		rgAPIv1.Use(s.httpHelpers.Middleware.BasicAuth(ctx, s.cfg.APIGW.APIServer.BasicAuth.Users))
	}

	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/upload", s.endpointUpload)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/notification", s.endpointNotification)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPut, "/document/identity", s.endpointAddDocumentIdentity)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodDelete, "/document/identity", s.endpointDeleteDocumentIdentity)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodDelete, "/document", s.endpointDeleteDocument)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/document/collect_id", s.endpointGetDocumentCollectID)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/identity/mapping", s.endpointIdentityMapping)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/document/list", s.endpointDocumentList)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/document", s.endpointGetDocument)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/document/search", s.endpointSearchDocuments)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/consent", s.endpointAddConsent)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/consent/get", s.endpointGetConsent)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/document/revoke", s.endpointRevokeDocument)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodPost, "/credential", s.endpointCredential)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIv1, http.MethodGet, "/credential/.well-known/jwks", s.endpointJWKS)

	// Run http server
	go func() {
		err := s.httpHelpers.Server.ListenAndServe(ctx, s.server, s.cfg.APIGW.APIServer)
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
