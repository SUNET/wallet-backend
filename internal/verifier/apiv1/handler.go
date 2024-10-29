package apiv1

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"vc/internal/gen/status/apiv1_status"
	"vc/pkg/model"
)

type JWTParts struct {
	Header    string
	Payload   string
	Signature string
}

type VerifyCredentialRequest struct {
	Jwt string `json:"jwt" validate:"required"`
}

type VerifyCredentialReply struct {
	Valid  bool   `json:"valid" validate:"required"`
	Reason string `json:"reason,omitempty"`
}

// Status return status for each ladok instance
func (c *Client) Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}
	status := probes.Check("verifier")
	return status, nil
}

func (c *Client) VerifyCredential(ctx context.Context, request *VerifyCredentialRequest) (*VerifyCredentialReply, error) {
	jwtParts, err := splitJWT(request.Jwt)
	if err != nil {
		msg := "Invalid JWT structure, unable to split into header, payload and signature"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Reason: msg}, nil
	}

	//c.log.Debug("jwt", "header", jwtParts.Header)
	//c.log.Debug("jwt", "payload", jwtParts.Payload)
	//c.log.Debug("jwt", "signature", jwtParts.Signature)

	fmt.Println("Header:", jwtParts.Header)
	fmt.Println("Payload:", jwtParts.Payload)
	fmt.Println("Signature:", jwtParts.Signature)

	//TODO(mk): impl more logic to verify jwt and set reply values

	return &VerifyCredentialReply{Valid: true}, nil
}

func splitJWT(jwt string) (*JWTParts, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT structure, expected three parts")
	}

	header, err := decodeBase64URL(parts[0])
	if err != nil {
		return nil, err
	}

	payload, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, err
	}

	signature := parts[2]

	return &JWTParts{
		Header:    header,
		Payload:   payload,
		Signature: signature,
	}, nil
}

func decodeBase64URL(encoded string) (string, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
