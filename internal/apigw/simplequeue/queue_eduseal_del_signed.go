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

// EduSealDelSigned holds the ladok delete signed queue
type EduSealDelSigned struct {
	service              *Service
	log                  *logger.Log
	metricEnqueueCounter prometheus.Counter
	*retask.Queue
}

// NewEduSealDelSigned creates a new ladok delete signed queue
func NewEduSealDelSigned(ctx context.Context, service *Service, queueName string, log *logger.Log) (*EduSealDelSigned, error) {
	eduSealDelSigned := &EduSealDelSigned{
		service: service,
		log:     log,
	}

	eduSealDelSigned.Queue = eduSealDelSigned.service.queueClient.NewQueue(ctx, queueName)

	eduSealDelSigned.metricEnqueueCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "apigw_queue_eduseal_del_signed_enqueue_total",
		Help: "The total number of added messages to the eduseal_del_signed queue",
	})

	eduSealDelSigned.log.Info("Started")

	return eduSealDelSigned, nil
}

// Enqueue publishes a document to the queue
func (s *EduSealDelSigned) Enqueue(ctx context.Context, message any) (*retask.Job, error) {
	s.log.Info("Enqueue delete signed pdf")
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealDelSigned:Enqueue")
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
func (s *EduSealDelSigned) Dequeue(ctx context.Context) error {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealDelSigned:Dequeue")
	defer span.End()
	return nil
}

// Wait waits for the next message
func (s *EduSealDelSigned) Wait(ctx context.Context) (*retask.Task, error) {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:EduSealDelSigned:Wait")
	defer span.End()

	task, err := s.Queue.Wait(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return task, nil
}
