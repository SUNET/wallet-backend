package apiv1

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"strings"
	"vc/internal/gen/status/apiv1_status"
	"vc/pkg/model"
)

// Status return status for each instance
func (c *Client) Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}
	status := probes.Check("verifier")
	return status, nil
}

type VerifyCredentialRequest struct {
	Jwt string `json:"jwt" validate:"required"`
}

type VerifyCredentialReply struct {
	Valid   bool   `json:"valid" validate:"required"`
	Message string `json:"message,omitempty"`
}

type JWTParts struct {
	Header    string
	Payload   string
	Signature string
}

type JWK struct {
	Crv string `json:"crv"`
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type PayloadCnf struct {
	Cnf struct {
		Jwk JWK `json:"jwk"`
	} `json:"cnf"`
}

func (c *Client) VerifyCredential(ctx context.Context, request *VerifyCredentialRequest) (*VerifyCredentialReply, error) {
	jwtParts, err := splitJWT(request.Jwt)
	if err != nil {
		msg := "Invalid JWT structure, unable to split into header, payload and signature"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Message: msg}, nil
	}

	//c.log.Debug("jwt", "header", jwtParts.Header)
	//c.log.Debug("jwt", "payload", jwtParts.Payload)
	//c.log.Debug("jwt", "signature", jwtParts.Signature)

	fmt.Println("Header:", jwtParts.Header)
	fmt.Println("Payload:", jwtParts.Payload)
	fmt.Println("Signature:", jwtParts.Signature)

	var payloadCnf PayloadCnf
	if err := json.Unmarshal([]byte(jwtParts.Payload), &payloadCnf); err != nil {
		msg := "Failed to parse payload JSON"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Message: msg}, nil
	}

	xStr := payloadCnf.Cnf.Jwk.X
	yStr := payloadCnf.Cnf.Jwk.Y

	fmt.Println("JWK x:", xStr)
	fmt.Println("JWK y:", yStr)

	xBytes, err := base64.RawURLEncoding.DecodeString(xStr)
	if err != nil {
		msg := "Error decoding x"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Message: msg}, nil
	}
	yBytes, err := base64.RawURLEncoding.DecodeString(yStr)
	if err != nil {
		msg := "Error decoding y"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Message: msg}, nil
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}

	token, err := jwt.Parse(request.Jwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		msg := "Error verifying token"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Message: msg}, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.log.Debug("Token is valid. Claims:")
		for key, val := range claims {
			c.log.Debug("claim", "key", key, "val", val)
		}
		return &VerifyCredentialReply{Valid: true}, nil
	}

	msg := "Invalid token"
	c.log.Debug(msg, "err", err)
	return &VerifyCredentialReply{Valid: false, Message: msg}, nil
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
