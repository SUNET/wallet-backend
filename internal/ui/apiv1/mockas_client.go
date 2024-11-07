package apiv1

import (
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"
)

type MockASClient struct {
	*VCBaseClient
}

func NewMockASClient(cfg *model.Cfg, tracer *trace.Tracer, logger *logger.Log) *MockASClient {
	return &MockASClient{
		VCBaseClient: NewClient("MOCKAS", cfg.UI.Services.MockAS.BaseURL, tracer, logger),
	}
}

func (c *MockASClient) Health() (any, error) {
	reply, err := c.DoGetJSON("/health")
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *MockASClient) MockNext(req *MockNextRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/mock/next", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
