package apiv1

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"os"
	"time"
	"vc/internal/gen/issuer/apiv1_issuer"
	"vc/pkg/helpers"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"
)

// Client holds the public api object
type Client struct {
	cfg            *model.Cfg
	tracer         *trace.Tracer
	log            *logger.Log
	apigwClient    *APIGWClient
	mockasClient   *MockASClient
	verifierClient *VerifierClient
	eventPublisher EventPublisher

	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	jwk        *apiv1_issuer.Jwk
}

// New creates a new instance of user interface web page
func New(ctx context.Context, cfg *model.Cfg, tracer *trace.Tracer, eventPublisher EventPublisher, log *logger.Log) (*Client, error) {
	c := &Client{
		cfg:            cfg,
		tracer:         tracer,
		log:            log.New("apiv1"),
		apigwClient:    NewAPIGWClient(cfg, tracer, log.New("apiwg_client")),
		mockasClient:   NewMockASClient(cfg, tracer, log.New("mockas_client")),
		verifierClient: NewVerifierClient(cfg, tracer, log.New("verifier_client")),
		eventPublisher: eventPublisher,
		jwk:            &apiv1_issuer.Jwk{},
	}

	err := c.initKeys(ctx)
	if err != nil {
		return nil, err
	}

	c.log.Info("Started")

	return c, nil
}

func (c *Client) initKeys(ctx context.Context) error {
	keyByte, err := os.ReadFile(c.cfg.Issuer.SigningKeyPath)
	if err != nil {
		c.log.Error(err, "Failed to read signing key, please create a ECDSA prime256v1 key and save it to the path")
		return err
	}

	if keyByte == nil {
		return helpers.ErrPrivateKeyMissing
	}

	c.privateKey, err = jwt.ParseECPrivateKeyFromPEM(keyByte)
	if err != nil {
		return err
	}

	c.publicKey = &c.privateKey.PublicKey

	if err := c.createJWK(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Client) createJWK(ctx context.Context) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(2*time.Second))
	defer cancel()

	key, err := jwk.New(c.privateKey)
	if err != nil {
		return err
	}

	key.Set("kid", "default_signing_key_id")

	if c.cfg.Issuer.JWTAttribute.Kid != "" {
		key.Set("kid", c.cfg.Issuer.JWTAttribute.Kid)
	}

	var jwkBytes []byte
	jwkBytes, err = json.MarshalIndent(key, "", "  ")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jwkBytes, c.jwk); err != nil {
		return err
	}

	return nil
}
