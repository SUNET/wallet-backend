package httpserver

import (
	"context"
	"fmt"
	"time"
	"vc/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Service) middlewareDuration(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		duration := time.Since(t)
		c.Set("duration", duration)
		c.Header("duration", duration.String())
	}
}

func (s *Service) middlewareTraceID(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("req_id", uuid.NewString())
		c.Header("req_id", c.GetString("req_id"))
		c.Next()
	}
}

func (s *Service) middlewareLogger(ctx context.Context) gin.HandlerFunc {
	log := s.logger.New("http")
	return func(c *gin.Context) {
		c.Next()
		log.Info("request", "status", c.Writer.Status(), "url", c.Request.URL.String(), "method", c.Request.Method, "req_id", c.GetString("req_id"))
	}
}

func (s *Service) middlewareAuthLog(ctx context.Context) gin.HandlerFunc {
	log := s.logger.New("http")
	return func(c *gin.Context) {
		u, _ := c.Get("user")
		c.Next()
		log.Info("auth", "user", u, "req_id", c.GetString("req_id"))
	}
}

func (s *Service) middlewareValidationCert(ctx context.Context) gin.HandlerFunc {
	log := s.logger.New("http")
	return func(c *gin.Context) {
		s.server.TLSConfig = s.tlsConfig
		c.Next()
		//log.Info("cert", "status", c.Writer.Status(), "url", c.Request.URL.String(), "method", c.Request.Method, "req_id", c.GetString("req_id"))
		log.Info("applying TLS config")
	}
}

func (s *Service) middlewareCrash(ctx context.Context) gin.HandlerFunc {
	log := s.logger.New("http")
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				status := c.Writer.Status()
				log.Trace("crash", "error", r, "status", status, "url", c.Request.URL.Path, "method", c.Request.Method)
				s.renderContent(ctx, c, 500, gin.H{"data": nil, "error": helpers.NewError("internal_server_error")})
			}
		}()
		c.Next()
	}
}

func (s *Service) middlewareClientCertAuth(ctx context.Context) gin.HandlerFunc {
	log := s.logger.New("http")
	return func(c *gin.Context) {
		clientCertSHA1 := c.Request.Header.Get("X-SSL-Client-SHA1")
		log.Info("clientCertSHA1", "clientCertSHA1", clientCertSHA1)
		fmt.Println("clientCertSHA1", clientCertSHA1)
		c.Next()
	}
}
