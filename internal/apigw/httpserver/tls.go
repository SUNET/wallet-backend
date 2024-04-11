package httpserver

import (
	"context"
	"crypto/tls"
)

func (s *Service) applyTLSConfig(ctx context.Context) {
	ctx, span := s.tp.Start(ctx, "httpserver:applyTLSConfig")
	defer span.End()

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	}

	s.server.TLSConfig = cfg
}
