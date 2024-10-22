package httpserver

import (
	"context"
	"vc/internal/gen/status/apiv1_status"
	"vc/internal/verifier/apiv1"
)

type Apiv1 interface {
	Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error)

	VerifyCredential(ctx context.Context, request *apiv1.VerifyCredentialRequest) (*apiv1.VerifyCredentialReply, error)
}
