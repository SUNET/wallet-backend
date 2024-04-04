package db

import (
	"context"
	"errors"
	"time"

	apiv1_status "vc/internal/gen/status/apiv1.status"
	"vc/pkg/logger"
	"vc/pkg/model"
	"vc/pkg/trace"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// ErrNoDocuments is returned when no documents are found
	ErrNoDocuments = errors.New("no documents in result")
)

// DB is the interface for the database
type DB interface {
	//connect(ctx context.Context) error
	//Save(ctx context.Context, doc *model.Document) error
	//Get(ctx context.Context, transactionID string) (*model.Document, error)
}

// Service is the database service
type Service struct {
	dbClient   *mongo.Client
	cfg        *model.Cfg
	log        *logger.Log
	tp         *trace.Tracer
	probeStore *apiv1_status.StatusProbeStore

	EduSealSigningColl *EduSealSigningColl
	VCDatastoreColl    *VCDatastoreColl
}

// New creates a new database service
func New(ctx context.Context, cfg *model.Cfg, tp *trace.Tracer, log *logger.Log) (*Service, error) {
	service := &Service{
		log:        log,
		cfg:        cfg,
		tp:         tp,
		probeStore: &apiv1_status.StatusProbeStore{},
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if err := service.connect(ctx); err != nil {
		return nil, err
	}

	var err error
	service.EduSealSigningColl, err = newEduSealSigningColl(ctx, service, service.dbClient.Database("eduseal").Collection("documents"))
	if err != nil {
		return nil, err
	}

	service.VCDatastoreColl, err = newVCDatastoreColl(ctx, service, service.dbClient.Database("vc").Collection("datastore"))
	if err != nil {
		return nil, err
	}

	service.log.Info("Started")

	return service, nil
}

// connect connects to the database
func (s *Service) connect(ctx context.Context) error {
	ctx, span := s.tp.Start(ctx, "apigw:db:connect")
	defer span.End()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(s.cfg.Common.Mongo.URI))
	if err != nil {
		return err
	}
	s.dbClient = client

	return nil
}

// Status returns the status of the database
func (s *Service) Status(ctx context.Context) *apiv1_status.StatusProbe {
	ctx, span := s.tp.Start(ctx, "db:status")
	defer span.End()

	if time.Now().Before(s.probeStore.NextCheck.AsTime()) {
		return s.probeStore.PreviousResult
	}
	probe := &apiv1_status.StatusProbe{
		Name:          "db",
		Healthy:       true,
		Message:       "OK",
		LastCheckedTS: timestamppb.Now(),
	}

	if err := s.dbClient.Ping(ctx, nil); err != nil {
		probe.Message = err.Error()
		probe.Healthy = false
	}

	s.probeStore.PreviousResult = probe
	s.probeStore.NextCheck = timestamppb.New(time.Now().Add(10 * time.Second))

	return probe
}

// Close closes the database connection
func (s *Service) Close(ctx context.Context) error {
	if err := s.dbClient.Disconnect(ctx); err != nil {
		return err
	}
	ctx.Done()
	return nil
}
