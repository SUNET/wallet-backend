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

const (
	ErrNotAJWT             = "not a jwt"
	ErrOnlySDJWTSupported  = "only typ sd-jwt or sd-jwt-vc is currently supported"
	ErrOnlyES256Supported  = "only alg ES256 is currently supported"
	ErrInvalidJWTStructure = "invalid JWT structure: expected format is header.payload.signature with optional ~disclosure~ segments (e.g., ~disclosure1~disclosure2~)."
	ErrInvalidPayloadJSON  = "failed to parse payload JSON"
	ErrInvalidJWKField     = "missing or invalid JWK field"
	ErrTokenVerification   = "error verifying token"
	ErrInvalidToken        = "invalid token"
)

type VerifyCredentialRequest struct {
	// min 20 is the ~ teoretical minimum for a non signed jwt encoded in base64
	Credential string `json:"credential" validate:"required,min=20"`
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

// VerifyCredential verifies a credential (only sd-jwt or sd-jwt-vc signed with ES256 is currently supported)
func (c *Client) VerifyCredential(ctx context.Context, request *VerifyCredentialRequest) (*VerifyCredentialReply, error) {
	jwtHeader, err := parseJWTHeader(request.Credential)
	if err != nil {
		return c.createInvalidReply(ErrNotAJWT, errors.New("Not a JWT"))
	}

	if !isTypSupported(jwtHeader) {
		return c.createInvalidReply(ErrOnlySDJWTSupported, errors.New("Typ not supported"))
	}

	if !isAlgSupported(jwtHeader) {
		return c.createInvalidReply(ErrOnlyES256Supported, errors.New("Alg not supported"))
	}

	sdjwtParts, err := splitSDJWT(request.Credential)
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
		c.debugLogClaims(claims)
		//TODO(mk): Verify that this verifier trusts the public key etc - as of now, the jwk-data used to create the pubkey is extracted from the jwt's payload.cnf.*, ie verify key binding/more info taken from the jwt.header etc!!!
		return &VerifyCredentialReply{Valid: true}, nil
	}

	return c.createInvalidReply(ErrInvalidToken, err)
}

type jwtHeader struct {
	Alg string `json:"alg" validate:"required,alg"`
	Typ string `json:"typ" validate:"required,typ"`
}

func parseJWTHeader(credential string) (*jwtHeader, error) {
	parts := strings.Split(credential, ".")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT header: %w", err)
	}

	var header jwtHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to parse JWT header: %w", err)
	}

	return &header, nil
}

func isTypSupported(header *jwtHeader) bool {
	return header.Typ == "sd-jwt" || header.Typ == "sd-jwt-vc"
}

func isAlgSupported(header *jwtHeader) bool {
	return header.Alg == "ES256"
}

func (c *Client) createInvalidReply(message string, err error) (*VerifyCredentialReply, error) {
	c.log.Debug(message, "err", err)
	return &VerifyCredentialReply{Valid: false, Message: message}, nil
}

func (c *Client) debugLogClaims(claims jwt.MapClaims) {
	c.log.Debug("Token is valid. Claims:")
	for key, val := range claims {
		c.log.Debug("claim", "key", key, "val", val)
	}
}

func splitSDJWT(credential string) (*SDJWTParts, error) {
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

func decodeBase64URL(encoded string) (string, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
