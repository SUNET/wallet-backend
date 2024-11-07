package httpserver

import (
	"context"
	"net/http"
	"vc/internal/ui/apiv1"
	"vc/pkg/httphelpers"
	"vc/pkg/trace"

	"vc/pkg/logger"
	"vc/pkg/model"

	"github.com/gin-gonic/gin"
)

// Service is the service object for httpserver
type Service struct {
	cfg           *model.Cfg
	log           *logger.Log
	tracer        *trace.Tracer
	server        *http.Server
	apiv1         Apiv1
	gin           *gin.Engine
	sessionConfig *sessionConfig
	httpHelpers   *httphelpers.Client
}

// sessionConfig... values is also used for the session cookie
type sessionConfig struct {
	//if name is changed, the web (javascript) must also be updated with the new name
	name                       string
	inactivityTimeoutInSeconds int
	path                       string
	httpOnly                   bool
	secure                     bool
	sameSite                   http.SameSite
	usernameKey                string
	loggedInTimeKey            string
}

// New creates a new httpserver service
func New(ctx context.Context, cfg *model.Cfg, apiv1 *apiv1.Client, tracer *trace.Tracer, log *logger.Log) (*Service, error) {
	s := &Service{
		cfg:    cfg,
		log:    log.New("httpserver"),
		tracer: tracer,
		apiv1:  apiv1,
		gin:    gin.New(),
		server: &http.Server{},
		sessionConfig: &sessionConfig{
			name:                       "vc_ui_auth_session",
			inactivityTimeoutInSeconds: cfg.UI.SessionInactivityTimeoutInSeconds,
			path:                       "/",
			httpOnly:                   true,
			secure:                     cfg.UI.APIServer.TLS.Enabled,
			sameSite:                   http.SameSiteStrictMode,
			usernameKey:                "username_key",
			loggedInTimeKey:            "logged_in_time_key",
		},
	}

	var err error
	s.httpHelpers, err = httphelpers.New(ctx, s.tracer, s.cfg, s.log)
	if err != nil {
		return nil, err
	}

	// extra middlewares (must be declared before Server.Default)
	s.gin.Use(s.httpHelpers.Middleware.Gzip(ctx))
	s.gin.Use(s.middlewareUserSession(ctx, s.cfg))

	rgRoot, err := s.httpHelpers.Server.Default(ctx, s.server, s.gin, s.cfg.UI.APIServer.Addr)
	if err != nil {
		return nil, err
	}

	s.gin.Static("/static", "./static")
	s.gin.LoadHTMLFiles("./static/index.html")
	s.gin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	s.httpHelpers.Server.RegEndpoint(ctx, rgRoot, http.MethodPost, "login", s.endpointLogin)
	s.httpHelpers.Server.RegEndpoint(ctx, rgRoot, http.MethodGet, "health", s.endpointHealth)

	rgAPIGW := rgRoot.Group("apigw")
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIGW, http.MethodGet, "health", s.endpointHealthAPIGW)

	rgVerifier := rgRoot.Group("verifier")
	s.httpHelpers.Server.RegEndpoint(ctx, rgVerifier, http.MethodGet, "health", s.endpointHealthVerifier)
	// No auth is needed to verify a credential
	s.httpHelpers.Server.RegEndpoint(ctx, rgVerifier, http.MethodPost, "verify", s.endpointVerifyCredential)

	rgMockAS := rgRoot.Group("mockas")
	s.httpHelpers.Server.RegEndpoint(ctx, rgMockAS, http.MethodGet, "health", s.endpointHealthMockAS)

	rgSecure := rgRoot.Group("secure", s.middlewareAuthRequired(ctx))
	s.httpHelpers.Server.RegEndpoint(ctx, rgSecure, http.MethodDelete, "logout", s.endpointLogout)
	s.httpHelpers.Server.RegEndpoint(ctx, rgSecure, http.MethodGet, "user", s.endpointUser)

	rgMockASSecure := rgSecure.Group("mockas")
	s.httpHelpers.Server.RegEndpoint(ctx, rgMockASSecure, http.MethodPost, "mock/next", s.endpointMockNext)

	rgAPIGWSecure := rgSecure.Group("apigw")
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIGWSecure, http.MethodPost, "document/list", s.endpointDocumentList)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIGWSecure, http.MethodPost, "upload", s.endpointUpload)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIGWSecure, http.MethodPost, "credential", s.endpointCredential)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIGWSecure, http.MethodPost, "document", s.endpointGetDocument)
	s.httpHelpers.Server.RegEndpoint(ctx, rgAPIGWSecure, http.MethodPost, "notification", s.endpointNotification)

	// Run http server
	go func() {
		err := s.httpHelpers.Server.ListenAndServe(ctx, s.server, s.cfg.UI.APIServer)
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
