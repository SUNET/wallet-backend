package httpserver

import (
	"context"
	apiv1_apigw "vc/internal/apigw/apiv1"
	"vc/internal/gen/status/apiv1_status"
	apiv1_mockas "vc/internal/mockas/apiv1"
	"vc/internal/ui/apiv1"
	apiv1_verifier "vc/internal/verifier/apiv1"
	"vc/pkg/model"
)

type Apiv1 interface {
	// ui
	Health(ctx context.Context, request *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error)
	Login(ctx context.Context, request *apiv1.LoginRequest) (*apiv1.LoggedinReply, error)
	Logout(ctx context.Context) error
	User(ctx context.Context) (*apiv1.LoggedinReply, error)

	// apigw
	HealthAPIGW(ctx context.Context, request *apiv1_status.StatusRequest) (any, error)
	DocumentList(ctx context.Context, request *apiv1.DocumentListRequest) (*apiv1_apigw.DocumentListReply, error)
	Upload(ctx context.Context, request *apiv1_apigw.UploadRequest) (any, error)
	Credential(ctx context.Context, request *apiv1.CredentialRequest) (any, error)
	GetDocument(ctx context.Context, request *apiv1.GetDocumentRequest) (any, error)
	Notification(ctx context.Context, reguest *apiv1.NotificationRequest) (any, error)
	SearchDocuments(ctx context.Context, request *model.SearchDocumentsRequest) (*model.SearchDocumentsReply, error)
	DeleteDocument(ctx context.Context, request *apiv1_apigw.DeleteDocumentRequest) error

	// mockas
	HealthMockAS(ctx context.Context, request *apiv1_status.StatusRequest) (any, error)
	MockNext(ctx context.Context, request *apiv1_mockas.MockNextRequest) (any, error)

	// verifier
	HealthVerifier(ctx context.Context, request *apiv1_status.StatusRequest) (any, error)
	Verify(ctx context.Context, request *apiv1_verifier.VerifyCredentialRequest) (any, error)
	DecodeCredential(ctx context.Context, request *apiv1_verifier.DecodeCredentialRequest) (any, error)
}
