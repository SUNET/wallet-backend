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

// EduSealValidate is the ladok unsigned queue
type EduSealValidate struct {
	service              *Service
	queueName            string
	log                  *logger.Log
	metricEnqueueCounter prometheus.Counter
	*retask.Queue
}

// NewEduSealValidate creates a new ladok unsigned queue
func NewEduSealValidate(ctx context.Context, service *Service, queueName string, log *logger.Log) (*EduSealValidate, error) {
	eduSealValidate := &EduSealValidate{
		service: service,
		log:     log,
	}

	eduSealValidate.Queue = eduSealValidate.service.queueClient.NewQueue(ctx, queueName)

	eduSealValidate.metricEnqueueCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_queue_eduseal_validate_enqueue_total",
		Help: "The total number of added messages to the eduseal_validate queue",
	})

	eduSealValidate.log.Info("Started")

	return eduSealValidate, nil
}

// Enqueue publishes a document to the queue
func (s *EduSealValidate) Enqueue(ctx context.Context, message any) (*retask.Job, error) {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealValidate:Enqueue")
	defer span.End()

	s.metricEnqueueCounter.Inc()

	data, err := json.Marshal(message)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return s.Queue.Enqueue(ctx, data)
}

// Wait waits for the next message
func (s *EduSealValidate) Wait(ctx context.Context) (*retask.Task, error) {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealValidate:Wait")
	defer span.End()

	return s.Queue.Wait(ctx)
}

// Dequeue dequeues a document from the queue
func (s *EduSealValidate) Dequeue(ctx context.Context) error {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealValidate:Dequeue")
	defer span.End()
	return nil
}
