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

// EduSealSign is the ladok unsigned queue
type EduSealSign struct {
	service              *Service
	log                  *logger.Log
	metricEnqueueCounter prometheus.Counter
	*retask.Queue
}

// NewEduSealSign creates a new ladok unsigned queue
func NewEduSealSign(ctx context.Context, service *Service, queueName string, log *logger.Log) (*EduSealSign, error) {
	eduSealSign := &EduSealSign{
		service: service,
		log:     log,
	}

	eduSealSign.Queue = eduSealSign.service.queueClient.NewQueue(ctx, queueName)

	eduSealSign.metricEnqueueCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_queue_eduseal_sign_enqueue_total",
		Help: "The total number of added messages to the eduseal_sign queue",
	})

	eduSealSign.log.Info("Started")

	return eduSealSign, nil
}

// Enqueue publishes a document to the queue
func (s *EduSealSign) Enqueue(ctx context.Context, message any) (*retask.Job, error) {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealSign:Enqueue")
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
func (s *EduSealSign) Dequeue(ctx context.Context) error {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealSign:Dequeue")
	defer span.End()
	return nil
}

// Wait waits for the next message
func (s *EduSealSign) Wait(ctx context.Context) (*retask.Task, error) {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealSign:Wait")
	defer span.End()
	return nil, nil
}
