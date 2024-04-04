package simplequeue

import (
	"context"
	"encoding/json"
	"vc/pkg/logger"

	retask "github.com/masv3971/goretask"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/codes"
)

// VCPersistentSave holds the ladok delete signed queue
type VCPersistentSave struct {
	service              *Service
	log                  *logger.Log
	metricEnqueueCounter prometheus.Counter
	*retask.Queue
}

// NewVCPersistentSave creates a new ladok delete signed queue
func NewVCPersistentSave(ctx context.Context, service *Service, queueName string, log *logger.Log) (*VCPersistentSave, error) {
	vcPersistentSave := &VCPersistentSave{
		service: service,
		log:     log,
	}

	vcPersistentSave.Queue = vcPersistentSave.service.queueClient.NewQueue(ctx, queueName)

	vcPersistentSave.metricEnqueueCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_vc_queue_persistent_save_enqueue_total",
		Help: "The total number of added messages to the persistent_save queue",
	})

	vcPersistentSave.log.Info("Started")

	return vcPersistentSave, nil
}

// Enqueue publishes a document to the queue
func (s *VCPersistentSave) Enqueue(ctx context.Context, message any) (*retask.Job, error) {
	s.log.Info("Enqueue")
	ctx, span := s.service.tp.Start(ctx, "simplequeue:VCPersistentSave:Enqueue")
	defer span.End()

	s.metricEnqueueCounter.Inc()

	data, err := json.Marshal(message)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return s.Queue.Enqueue(ctx, data)
}

// Dequeue dequeues a document from the queue
func (s *VCPersistentSave) Dequeue(ctx context.Context) error {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:VCPersistentSave:Dequeue")
	defer span.End()
	return nil
}

// Wait waits for the next message
func (s *VCPersistentSave) Wait(ctx context.Context) (*retask.Task, error) {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:VCPersistentSave:Wait")
	defer span.End()

	task, err := s.Queue.Wait(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return task, nil
}
