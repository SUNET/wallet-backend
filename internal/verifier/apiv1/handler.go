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
	"time"
	"vc/internal/gen/status/apiv1_status"
	"vc/pkg/model"
)

// Status return status for each instance
func (c *Client) Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}
	return probes.Check("verifier"), nil
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

	c.log.Debug("jwt", "header", jwtParts.Header)
	c.log.Debug("jwt", "payload", jwtParts.Payload)
	c.log.Debug("jwt", "signature", jwtParts.Signature)

	//fmt.Println("Header:", jwtParts.Header)
	//fmt.Println("Payload:", jwtParts.Payload)
	//fmt.Println("Signature:", jwtParts.Signature)

	var payloadCnf PayloadCnf
	if err := json.Unmarshal([]byte(jwtParts.Payload), &payloadCnf); err != nil {
		msg := "Failed to parse payload JSON"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Message: msg}, nil
	}

	xStr := payloadCnf.Cnf.Jwk.X
	yStr := payloadCnf.Cnf.Jwk.Y

	c.log.Debug("jwt", "x", xStr, "y", yStr)

	//fmt.Println("JWK x:", xStr)
	//fmt.Println("JWK y:", yStr)

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

	token, err := jwt.Parse(request.Jwt, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	c.log.Debug("jwt", "token", token)

	if err != nil {
		msg := "Error verifying token"
		c.log.Debug(msg, "err", err)
		return &VerifyCredentialReply{Valid: false, Message: msg}, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := validateClaims(claims); err != nil {
			msg := "Invalid claims"
			c.log.Debug(msg, "err", err)
			return &VerifyCredentialReply{Valid: false, Message: msg}, nil
		}

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

func validateClaims(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return errors.New("token has expired")
		}
	}

	if nbf, ok := claims["nbf"].(float64); ok {
		if time.Now().Unix() < int64(nbf) {
			return errors.New("token is not valid yet")
		}
	}

	//if iss, ok := claims["iss"].(string); ok {
	//	if iss != "https://issuer.sunet.se" { //TODO: where do I find list of trusted issuers?
	//		return fmt.Errorf("unexpected issuer: %s", iss)
	//	}
	//}

	return nil
}

func decodeBase64URL(encoded string) (string, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
