package simplequeue

import (
	"context"
	"vc/pkg/kvclient"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"

	retask "github.com/masv3971/goretask"
	"github.com/redis/go-redis/v9"
)

type queue interface {
	Enqueue(ctx context.Context, message any) (*retask.Job, error)
	Dequeue(ctx context.Context) error
	Wait(ctx context.Context) (*retask.Task, error)
}

// Service is the service object for queue
type Service struct {
	redisClient *redis.Client
	queueClient *retask.Client
	tp          *trace.Tracer
	log         *logger.Log
	cfg         *model.Cfg

	EduSealSign           queue
	EduSealValidate       queue
	EduSealDelSigned      queue
	EduSealPersistentSave queue
	VCPersistentSave      queue
	VCPersistentGet       queue
	VCPersistentDelete    queue
}

// New creates a new queue service
func New(ctx context.Context, kv *kvclient.Client, tracer *trace.Tracer, cfg *model.Cfg, log *logger.Log) (*Service, error) {
	service := &Service{
		redisClient: kv.RedisClient,
		tp:          tracer,
		log:         log,
		cfg:         cfg,
	}
	var err error
	service.queueClient, err = retask.New(ctx, service.redisClient)
	if err != nil {
		return nil, err
	}

	service.EduSealSign, err = NewEduSealSign(ctx, service, cfg.Common.Queues.SimpleQueue.EduSealSign.Name, service.log.New("LadokSign"))
	if err != nil {
		return nil, err
	}

	service.EduSealValidate, err = NewEduSealValidate(ctx, service, cfg.Common.Queues.SimpleQueue.EduSealValidate.Name, service.log.New("LadokValidate"))
	if err != nil {
		return nil, err
	}

	service.EduSealDelSigned, err = NewEduSealDelSigned(ctx, service, cfg.Common.Queues.SimpleQueue.EduSealDelSigned.Name, service.log.New("LadokDelSigned"))
	if err != nil {
		return nil, err
	}

	service.EduSealPersistentSave, err = NewEduSealPersistentSave(ctx, service, cfg.Common.Queues.SimpleQueue.EduSealPersistentSave.Name, service.log.New("LadokPersistentSave"))
	if err != nil {
		return nil, err
	}

	service.VCPersistentSave, err = NewVCPersistentSave(ctx, service, cfg.Common.Queues.SimpleQueue.VCPersistentSave.Name, service.log.New("VCPersistentSave"))
	if err != nil {
		return nil, err
	}

	service.VCPersistentGet, err = NewVCPersistentGet(ctx, service, cfg.Common.Queues.SimpleQueue.VCPersistentGet.Name, service.log.New("VCPersistentGet"))
	if err != nil {
		return nil, err
	}

	service.VCPersistentDelete, err = NewVCPersistentDelete(ctx, service, cfg.Common.Queues.SimpleQueue.VCPersistentDelete.Name, service.log.New("VCPersistentDelete"))
	if err != nil {
		return nil, err
	}

	service.log.Info("Started")

	return service, nil
}

// Close closes the service
func (s *Service) Close(ctx context.Context) error {
	s.log.Info("Stopped")
	return nil
}
