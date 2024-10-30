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

const (
	ErrInvalidJWTStructure = "invalid JWT structure, expected header.payload.signature"
	ErrInvalidPayloadJSON  = "failed to parse payload JSON"
	ErrInvalidJWKField     = "missing or invalid JWK field"
	ErrTokenVerification   = "error verifying token"
	ErrInvalidToken        = "invalid token"
	ErrInvalidClaims       = "invalid claims"
	ErrExpiredToken        = "token has expired"
	ErrNotYetValidToken    = "token is not valid yet"
)

type VerifyCredentialRequest struct {
	Credential string `json:"credential" validate:"required"`
}

type VerifyCredentialReply struct {
	Valid   bool   `json:"valid" validate:"required"`
	Message string `json:"message,omitempty"`
}

type SDJWTParts struct {
	CompleteJWT string
	Header      string
	Payload     string
	Signature   string
	Disclosures []string
}

type JWK struct {
	Crv string `json:"crv"`
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

// Status returns the status for each instance.
func (c *Client) Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}
	return probes.Check("verifier"), nil
}

// VerifyCredential verifies a credential (sd-jwt currently supported)
func (c *Client) VerifyCredential(ctx context.Context, request *VerifyCredentialRequest) (*VerifyCredentialReply, error) {
	//TODO(mk): ev. avgör snabbt om det är en sd-jwt eller någon annan typ av credential?

	sdjwtParts, err := splitCredential(request.Credential)
	if err != nil {
		return c.createInvalidReply(ErrInvalidJWTStructure, err)
	}
	c.log.Debug("credential", "parts", sdjwtParts)

	jwk, err := extractJWK(sdjwtParts.Payload)
	if err != nil {
		return c.createInvalidReply(ErrInvalidJWKField, err)
	}
	c.log.Debug("jwk", "data", jwk)

	pubKey, err := createPubKey(jwk)
	if err != nil {
		return c.createInvalidReply(ErrInvalidJWKField, err)
	}
	c.log.Debug("pubkey", "data", pubKey)

	token, err := parseJWT(sdjwtParts.CompleteJWT, pubKey)
	if err != nil {
		return c.createInvalidReply(ErrTokenVerification, err)
	}
	c.log.Debug("token", "data", token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := validateClaims(claims); err != nil {
			return c.createInvalidReply(ErrInvalidClaims, err)
		}
		c.logClaims(claims)
		//TODO(mk): validate disclosures in any way?
		return &VerifyCredentialReply{Valid: true}, nil
	}

	return c.createInvalidReply(ErrInvalidToken, err)
}

func (c *Client) createInvalidReply(message string, err error) (*VerifyCredentialReply, error) {
	c.log.Debug(message, "err", err)
	return &VerifyCredentialReply{Valid: false, Message: message}, nil
}

func (c *Client) logClaims(claims jwt.MapClaims) {
	c.log.Debug("Token is valid. Claims:")
	for key, val := range claims {
		c.log.Debug("claim", "key", key, "val", val)
	}
}

func splitCredential(credential string) (*SDJWTParts, error) {
	parts := strings.Split(credential, "~")
	jwtParts := strings.Split(parts[0], ".")
	if len(jwtParts) != 3 {
		return nil, errors.New(ErrInvalidJWTStructure)
	}

	header, err := decodeBase64URL(jwtParts[0])
	if err != nil {
		return nil, err
	}

	payload, err := decodeBase64URL(jwtParts[1])
	if err != nil {
		return nil, err
	}

	return &SDJWTParts{
		CompleteJWT: parts[0],
		Header:      header,
		Payload:     payload,
		Signature:   jwtParts[2],
		Disclosures: parts[1:],
	}, nil
}

func extractJWK(payload string) (*JWK, error) {
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
		return nil, errors.New(ErrInvalidPayloadJSON)
	}

	cnf, ok := payloadMap["cnf"].(map[string]interface{})
	if !ok {
		return nil, errors.New(ErrInvalidJWKField)
	}

	jwkMap, ok := cnf["jwk"].(map[string]interface{})
	if !ok {
		return nil, errors.New(ErrInvalidJWKField)
	}

	jwk := &JWK{}
	if x, ok := jwkMap["x"].(string); ok {
		jwk.X = x
	} else {
		return nil, errors.New("missing or invalid 'x' field in JWK")
	}

	if y, ok := jwkMap["y"].(string); ok {
		jwk.Y = y
	} else {
		return nil, errors.New("missing or invalid 'y' field in JWK")
	}

	if crv, ok := jwkMap["crv"].(string); ok {
		jwk.Crv = crv
	} else {
		return nil, errors.New("missing or invalid 'crv' field in JWK")
	}

	if kty, ok := jwkMap["kty"].(string); ok {
		jwk.Kty = kty
	} else {
		return nil, errors.New("missing or invalid 'kty' field in JWK")
	}

	return jwk, nil
}

func createPubKey(jwk *JWK) (*ecdsa.PublicKey, error) {
	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("error decoding x: %w", err)
	}

	yBytes, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, fmt.Errorf("error decoding y: %w", err)
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}

	return pubKey, nil
}

func parseJWT(completeJWT string, pubKey *ecdsa.PublicKey) (*jwt.Token, error) {
	return jwt.Parse(completeJWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
}

func validateClaims(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return errors.New(ErrExpiredToken)
		}
	}

	if nbf, ok := claims["nbf"].(float64); ok {
		if time.Now().Unix() < int64(nbf) {
			return errors.New(ErrNotYetValidToken)
		}
	}

	//TODO(mk): validate more claims here ...
	return nil
}

func decodeBase64URL(encoded string) (string, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
