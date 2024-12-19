package apiv1

import (
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"
)

// APIGWClient client type for apigw
type APIGWClient struct {
	*VCBaseClient
}

// NewAPIGWClient creates a new api client for apigw
func NewAPIGWClient(cfg *model.Cfg, tracer *trace.Tracer, logger *logger.Log) *APIGWClient {
	return &APIGWClient{
		VCBaseClient: NewClient("APIGW", cfg.UI.Services.APIGW.BaseURL, tracer, logger),
	}
}

// SearchDocuments sends POST to /api/v1/document/search
func (c *APIGWClient) SearchDocuments(req *model.SearchDocumentsRequest) (*model.SearchDocumentsReply, error) {
	reply, err := DoPostJSONGeneric[model.SearchDocumentsReply](c.VCBaseClient, "/api/v1/document/search", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
