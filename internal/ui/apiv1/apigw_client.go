package apiv1

import (
	apiv1_apigw "vc/internal/apigw/apiv1"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"
)

type APIGWClient struct {
	*VCBaseClient
}

func NewAPIGWClient(cfg *model.Cfg, tracer *trace.Tracer, logger *logger.Log) *APIGWClient {
	return &APIGWClient{
		VCBaseClient: NewClient("APIGW", cfg.UI.Services.APIGW.BaseURL, tracer, logger),
	}
}

func (c *APIGWClient) DocumentList(req *DocumentListRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/document/list", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *APIGWClient) Status() (any, error) {
	reply, err := c.DoGetJSON("/health")
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *APIGWClient) Upload(req *apiv1_apigw.UploadRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/upload", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *APIGWClient) Credential(req *CredentialRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/credential", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *APIGWClient) GetDocument(req *GetDocumentRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/document", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *APIGWClient) Notification(request *NotificationRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/notification", request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
