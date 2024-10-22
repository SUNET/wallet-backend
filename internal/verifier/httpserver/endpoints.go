package httpserver

import (
	"context"
	"go.opentelemetry.io/otel/codes"
	"vc/internal/verifier/apiv1"

	"vc/internal/gen/status/apiv1_status"

	"github.com/gin-gonic/gin"
)

func (s *Service) endpointHealth(ctx context.Context, c *gin.Context) (any, error) {
	//TODO:(mk): fix endpointHealth
	request := &apiv1_status.StatusRequest{}
	reply, err := s.apiv1.Status(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (s *Service) endpointVerifyCredential(ctx context.Context, c *gin.Context) (any, error) {
	ctx, span := s.tracer.Start(ctx, "httpserver:endpointVerifyCredential")
	defer span.End()

	request := &apiv1.VerifyCredentialRequest{}
	if err := s.httpHelpers.Binding.Request(ctx, c, request); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	reply, err := s.apiv1.VerifyCredential(ctx, request)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return reply, nil
}
