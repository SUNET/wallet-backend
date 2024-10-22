package apiv1

import (
	"context"
	"vc/internal/gen/status/apiv1_status"
	"vc/pkg/model"
)

type VerifyCredentialRequest struct {
	//TODO(mk): komplett Credential
}

type VerifyCredentialReply struct {
}

// Status return status for each ladok instance
func (c *Client) Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}
	status := probes.Check("verifier")
	return status, nil
}

func (c *Client) VerifyCredential(ctx context.Context, request *VerifyCredentialRequest) (*VerifyCredentialReply, error) {
	reply := &VerifyCredentialReply{}
	//TODO(mk): impl logic for VerifyCredential
	return reply, nil
}
