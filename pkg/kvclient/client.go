package kvclient

import (
	"context"
	"time"
	apiv1_status "vc/internal/gen/status/apiv1.status"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/redis/go-redis/v9"
)

// Client holds the kv object
type Client struct {
	RedisClient *redis.Client
	cfg         *model.Cfg
	log         *logger.Log
	probeStore  *apiv1_status.StatusProbeStore
	tp          *trace.Tracer

	Doc *Doc
}

// New creates a new instance of kv
func New(ctx context.Context, cfg *model.Cfg, tracer *trace.Tracer, log *logger.Log) (*Client, error) {
	c := &Client{
		cfg:        cfg,
		log:        log,
		probeStore: &apiv1_status.StatusProbeStore{},
		tp:         tracer,
	}

	c.RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Common.KeyValue.Addr,
		Password: cfg.Common.KeyValue.Password,
		DB:       cfg.Common.KeyValue.DB,
	},
	)

	c.Doc = &Doc{
		client: c,
		key:    "doc:%s:%s",
	}

	c.log.Info("Started")

	return c, nil
}

// Status returns the status of the database
func (c *Client) Status(ctx context.Context) *apiv1_status.StatusProbe {
	ctx, span := c.tp.Start(ctx, "kvclient:Status")
	defer span.End()

	if time.Now().Before(c.probeStore.NextCheck.AsTime()) {
		return c.probeStore.PreviousResult
	}
	probe := &apiv1_status.StatusProbe{
		Name:          "kv",
		Healthy:       true,
		Message:       "OK",
		LastCheckedTS: timestamppb.Now(),
	}

	_, err := c.RedisClient.Ping(ctx).Result()
	if err != nil {
		probe.Message = err.Error()
		probe.Healthy = false
	}
	c.probeStore.PreviousResult = probe
	c.probeStore.NextCheck = timestamppb.New(time.Now().Add(time.Second * 10))

	return probe
}

// Close closes the connection to the database
func (c *Client) Close(ctx context.Context) error {
	return c.RedisClient.Close()
}
