package apiv1

import (
	"context"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"
)

// Client holds the public api object
type Client struct {
	cfg    *model.Cfg
	log    *logger.Log
	tracer *trace.Tracer

	apigwClient *APIGWClient
}

// New creates a new instance of the public api
func New(ctx context.Context, tracer *trace.Tracer, cfg *model.Cfg, log *logger.Log) (*Client, error) {
	c := &Client{
		cfg:         cfg,
		tracer:      tracer,
		log:         log.New("apiv1"),
		apigwClient: NewAPIGWClient(cfg, tracer, log.New("apiwg_client")),
	}

	c.log.Info("Started")

	return c, nil
}
