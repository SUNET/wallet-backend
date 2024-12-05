package httpserver

import (
	"context"
	apiv1_apigw "vc/internal/apigw/apiv1"
	"vc/internal/gen/status/apiv1_status"

	"github.com/gin-gonic/gin"
)

func (s *Service) endpointHealth(ctx context.Context, c *gin.Context) (any, error) {
	request := &apiv1_status.StatusRequest{}
	reply, err := s.apiv1.Status(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (s *Service) endpointSearchDocuments(ctx context.Context, c *gin.Context) (any, error) {
	request := &apiv1_apigw.SearchDocumentsRequest{}
	if err := s.httpHelpers.Binding.Request(ctx, c, request); err != nil {
		return nil, err
	}

	reply, err := s.apiv1.SearchDocuments(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
