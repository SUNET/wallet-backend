package apiv1

import (
	"context"
	"vc/internal/apigw/db"
	"vc/internal/apigw/simplequeue"
	"vc/pkg/kvclient"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"
)

//	@title		Datastore API
//	@version	0.1.0
//	@BasePath	/api/v1

// Client holds the public api object
type Client struct {
	cfg         *model.Cfg
	db          *db.Service
	log         *logger.Log
	tp          *trace.Tracer
	kv          *kvclient.Client
	simpleQueue *simplequeue.Service
}

// New creates a new instance of the public api
func New(ctx context.Context, kv *kvclient.Client, db *db.Service, simplequeue *simplequeue.Service, tp *trace.Tracer, cfg *model.Cfg, logger *logger.Log) (*Client, error) {
	c := &Client{
		cfg:         cfg,
		db:          db,
		log:         logger,
		kv:          kv,
		tp:          tp,
		simpleQueue: simplequeue,
	}

	c.log.Info("Started")

	return c, nil
}
