package apiv1

import (
	"context"
	apiv1_registry "vc/internal/gen/registry/apiv1.registry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CredentialRequest is the request for Credential
type CredentialRequest struct {
	AuthenticSource string `json:"authentic_source" binding:"required"`
	DocumentID      string `json:"document_id" binding:"required"`
	DocumentType    string `json:"document_type" binding:"required"`
	CollectID       string `json:"collect_id" binding:"required"`
	DateOfBirth     string `json:"date_of_birth" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
}

// CredentialReply is the reply for Credential
type CredentialReply struct {
	SDJWT string `json:"sdjwt"`
}

// Credential is the credential endpoint
func (c *Client) Credential(ctx context.Context, req *CredentialRequest) (*CredentialReply, error) {
	c.log.Info("Credential", "req", req)
	// IDMapping

	// GetDocument

	// Build SDJWT
	reply := &CredentialReply{
		SDJWT: "mock sd-jwt",
	}
	return reply, nil
}

// RevokeRequest is the request for GenericRevoke
type RevokeRequest struct {
	AuthenticSource string `json:"authentic_source"`
	DocumentType    string `json:"document_type"`
	DocumentID      string `json:"document_id"`
	RevocationID    string `json:"revocation_id"`
}

// RevokeReply is the reply for GenericRevoke
type RevokeReply struct {
	Data struct {
		Status bool `json:"status"`
	}
}

// Revoke revokes a document
//
//	@Summary		Revoke
//	@ID				generic-revoke
//	@Description	Revoke endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	RevokeReply				"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		RevokeRequest			true	" "
//	@Router			/revoke [post]
func (c *Client) Revoke(ctx context.Context, req *RevokeRequest) (*RevokeReply, error) {
	optInsecure := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(c.cfg.Registry.RPCServer.Addr, optInsecure)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := apiv1_registry.NewRegistryServiceClient(conn)
	resp, err := client.Revoke(ctx, &apiv1_registry.RevokeRequest{
		Entity: "mura",
	})
	if err != nil {
		return nil, err
	}

	reply := &RevokeReply{
		Data: struct {
			Status bool `json:"status"`
		}{
			Status: resp.Status,
		},
	}
	return reply, nil
}
