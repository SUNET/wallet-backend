package apiv1

import (
	"context"
	"errors"
	"time"
	apiv1_apigw "vc/internal/apigw/apiv1"
	"vc/internal/gen/status/apiv1_status"
	apiv1_mockas "vc/internal/mockas/apiv1"
	apiv1_verifier "vc/internal/verifier/apiv1"
	"vc/pkg/model"
)

func (c *Client) Health(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}
	status := probes.Check("ui")
	return status, nil
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoggedinReply struct {
	Username string `json:"username" validate:"required"`
	// LoggedInTime RFC3339
	LoggedInTime time.Time `json:"logged_in_time" validate:"required"`
}

func (c *Client) Login(ctx context.Context, req *LoginRequest) (*LoggedinReply, error) {
	if req.Username != c.cfg.UI.Username || req.Password != c.cfg.UI.Password {
		return nil, errors.New("invalid username and/or password")
	}

	reply := &LoggedinReply{
		Username:     c.cfg.UI.Username,
		LoggedInTime: time.Now(),
	}

	return reply, nil
}

func (c *Client) Logout(ctx context.Context) error {
	return nil
}

func (c *Client) User(ctx context.Context) (*LoggedinReply, error) {
	return nil, nil
}

type DocumentListRequest struct {
	AuthenticSource string          `json:"authentic_source"`
	Identity        *model.Identity `json:"identity" validate:"required"`
	DocumentType    string          `json:"document_type"`
	ValidFrom       int64           `json:"valid_from"`
	ValidTo         int64           `json:"valid_to"`
}

func (c *Client) DocumentList(ctx context.Context, req *DocumentListRequest) (*apiv1_apigw.DocumentListReply, error) {
	reply, err := c.apigwClient.DocumentList(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) Upload(ctx context.Context, req *apiv1_apigw.UploadRequest) (any, error) {
	reply, err := c.apigwClient.Upload(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// CredentialRequest is the request for the Credential endpoint
type CredentialRequest struct {
	AuthenticSource string          `json:"authentic_source" validate:"required"`
	Identity        *model.Identity `json:"identity" validate:"required"`
	DocumentType    string          `json:"document_type" validate:"required"`
	CredentialType  string          `json:"credential_type" validate:"required"`
	CollectID       string          `json:"collect_id" validate:"required"`
	JWK             map[string]any  `json:"jwk"`
}

// Credential sends POST to apigw /api/v1/credential
func (c *Client) Credential(ctx context.Context, req *CredentialRequest) (any, error) {
	req.JWK = map[string]any{
		"kty": c.jwk.Kty,
		"kid": c.jwk.Kid,
		"crv": c.jwk.Crv,
		"x":   c.jwk.X,
		"y":   c.jwk.Y,
	}

	reply, err := c.apigwClient.Credential(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GetDocumentRequest is the request for the GetDocument endpoint
type GetDocumentRequest struct {
	AuthenticSource string `json:"authentic_source" validate:"required"`
	DocumentType    string `json:"document_type" validate:"required"`
	DocumentID      string `json:"document_id" validate:"required"`
}

func (c *Client) GetDocument(ctx context.Context, req *GetDocumentRequest) (any, error) {
	reply, err := c.apigwClient.GetDocument(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

type NotificationRequest struct {
	AuthenticSource string `json:"authentic_source" validate:"required"`
	DocumentType    string `json:"document_type" validate:"required"`
	DocumentID      string `json:"document_id" validate:"required"`
}

func (c *Client) Notification(ctx context.Context, req *NotificationRequest) (any, error) {
	reply, err := c.apigwClient.Notification(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

//type MockNextRequest struct {
//	DocumentType            string `json:"document_type" validate:"required"`
//	AuthenticSource         string `json:"authentic_source" validate:"required"`
//	AuthenticSourcePersonId string `json:"authentic_source_person_id" validate:"required"`
//	IdentitySchemaName      string `json:"identity_schema_name" validate:"required"`
//}

func (c *Client) MockNext(ctx context.Context, req *apiv1_mockas.MockNextRequest) (any, error) {
	if c.cfg.Common.Kafka.Enabled {
		if err := c.eventPublisher.MockNext(req); err != nil {
			return nil, err
		}
		return nil, nil
	}

	reply, err := c.mockasClient.MockNext(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) Verify(ctx context.Context, req *apiv1_verifier.VerifyCredentialRequest) (any, error) {
	reply, err := c.verifierClient.Verify(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) HealthAPIGW(ctx context.Context, req *apiv1_status.StatusRequest) (any, error) {
	reply, err := c.apigwClient.Health()
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) HealthVerifier(ctx context.Context, req *apiv1_status.StatusRequest) (any, error) {
	reply, err := c.verifierClient.Health()
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) HealthMockAS(ctx context.Context, req *apiv1_status.StatusRequest) (any, error) {
	reply, err := c.mockasClient.Health()
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) DecodeCredential(ctx context.Context, req *apiv1_verifier.DecodeCredentialRequest) (any, error) {
	reply, err := c.verifierClient.Decode(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) SearchDocuments(ctx context.Context, req *apiv1_apigw.SearchDocumentsRequest) (*apiv1_apigw.SearchDocumentsReply, error) {
	reply, err := c.apigwClient.SearchDocuments(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *Client) DeleteDocument(ctx context.Context, req *apiv1_apigw.DeleteDocumentRequest) error {
	err := c.apigwClient.DeleteDocument(req)
	if err != nil {
		return err
	}
	return nil
}
