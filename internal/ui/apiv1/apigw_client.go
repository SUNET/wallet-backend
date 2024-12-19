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

// DocumentList sends POST to apigw /api/v1/document/list
func (c *APIGWClient) DocumentList(req *DocumentListRequest) (*apiv1_apigw.DocumentListReply, error) {
	reply, err := DoPostJSONGeneric[apiv1_apigw.DocumentListReply](c.VCBaseClient, "/api/v1/document/list", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Status sends GET to apigw /health
func (c *APIGWClient) Health() (any, error) {
	reply, err := c.DoGetJSON("/health")
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Upload sends POST to apigw /api/v1/upload
func (c *APIGWClient) Upload(req *apiv1_apigw.UploadRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/upload", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Credential sends POST to apigw /api/v1/credential
func (c *APIGWClient) Credential(req *CredentialRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/credential", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GetDocument sends POST to apigw /api/v1/document
func (c *APIGWClient) GetDocument(req *GetDocumentRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/document", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Notification sends POST to /api/v1/notification
func (c *APIGWClient) Notification(req *NotificationRequest) (any, error) {
	reply, err := c.DoPostJSON("/api/v1/notification", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// SearchDocuments sends POST to /api/v1/document/search
func (c *APIGWClient) SearchDocuments(req *model.SearchDocumentsRequest) (*model.SearchDocumentsReply, error) {
	reply, err := DoPostJSONGeneric[model.SearchDocumentsReply](c.VCBaseClient, "/api/v1/document/search", req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// DeleteDocument sends DELETE to /api/v1/document
func (c *APIGWClient) DeleteDocument(req *apiv1_apigw.DeleteDocumentRequest) error {
	err := c.DoDelete("/api/v1/document", req)
	if err != nil {
		return err
	}
	return nil
}
