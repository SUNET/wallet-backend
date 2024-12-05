package apiv1

import (
	apiv1_apigw "vc/internal/apigw/apiv1"
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

func (c *APIGWClient) SearchDocuments(req *apiv1_apigw.SearchDocumentsRequest) (*apiv1_apigw.SearchDocumentsReply, error) {
	reply, err := DoPostJSONGeneric[apiv1_apigw.SearchDocumentsReply](c.VCBaseClient, "/api/v1/document/search", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
