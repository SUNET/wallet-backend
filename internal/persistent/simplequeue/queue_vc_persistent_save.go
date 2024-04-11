package simplequeue

import (
	"context"
	"encoding/json"
	"vc/pkg/logger"
	"vc/pkg/model"

	retask "github.com/masv3971/goretask"
	"go.opentelemetry.io/otel/codes"
)

// VCPersistentSave holds the ladok delete signed queue
type VCPersistentSave struct {
	service *Service
	log     *logger.Log
	*retask.Queue
}

// NewVCPersistentSave creates a new ladok delete signed queue
func NewVCPersistentSave(ctx context.Context, service *Service, queueName string, log *logger.Log) (*VCPersistentSave, error) {
	vcPersistentSave := &VCPersistentSave{
		service: service,
		log:     log,
	}

	vcPersistentSave.Queue = vcPersistentSave.service.queueClient.NewQueue(ctx, queueName)

	vcPersistentSave.log.Info("Started")

	return vcPersistentSave, nil
}

// Enqueue publishes a document to the queue
func (s *VCPersistentSave) Enqueue(ctx context.Context, message any) (*retask.Job, error) {
	s.log.Info("Enqueue")
	ctx, span := s.service.tp.Start(ctx, "simplequeue:VCPersistentSave:Enqueue")
	defer span.End()

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

// Worker is the worker
func (s *VCPersistentSave) Worker(ctx context.Context) error {
	ctx, span := s.service.tp.Start(ctx, "simplequeue:VCPersistentSave:Worker")
	defer span.End()

	var (
		taskChan = make(chan *retask.Task)
		errChan  = make(chan error)
	)

	go func() {
		for {
			task, err := s.Wait(ctx)
			if err != nil {
				errChan <- err
			}
			taskChan <- task
		}
	}()

	for {
		select {
		case err := <-errChan:
			s.log.Error(err, "Worker failed")
			return err
		case task := <-taskChan:
			s.log.Info("Got task", "task", task.Data)
			document := &model.Upload{}
			if err := json.Unmarshal([]byte(task.Data), document); err != nil {
				span.SetStatus(codes.Error, err.Error())
				s.log.Error(err, "Unmarshal failed")
			}
			if err := s.service.db.VCDatastoreColl.Save(ctx, document); err != nil {
				span.SetStatus(codes.Error, err.Error())
				s.log.Error(err, "Save upload failed")
			}

		case <-ctx.Done():
			s.log.Info("Stopped worker")
			return nil
		}
	}
}
