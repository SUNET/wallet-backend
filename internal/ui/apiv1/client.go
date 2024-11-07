package apiv1

import (
	"context"
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
	}

	c.log.Info("Started")

	return c, nil
}
