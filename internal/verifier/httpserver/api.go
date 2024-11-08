package httpserver

import (
	"context"
	"vc/internal/gen/status/apiv1_status"
	"vc/internal/verifier/apiv1"
)

type Apiv1 interface {
	Health(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error)
	VerifyCredential(ctx context.Context, request *apiv1.Credential) (*apiv1.VerifyCredentialReply, error)
	DecodeCredential(ctx context.Context, request *apiv1.Credential) (*apiv1.DecodedCredential, error)
}
