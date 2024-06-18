package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"time"
)

type Server struct {
	renewInterval time.Duration
	connectorRepo repository.ConnectorRepository
	docRepo       repository.DocumentRepository
	messenger     messaging.Client
	scheduler     gocron.Scheduler
	streamCfg     *messaging.StreamConfig
	fileSizeLimit int
}

func NewServer(
	cfg *Config,
	connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	messenger messaging.Client,
	messagingCfg *messaging.Config) (*Server, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Server{connectorRepo: connectorRepo,
		docRepo:       docRepo,
		renewInterval: time.Duration(cfg.RenewInterval) * time.Second,
		fileSizeLimit: cfg.FileSizeLimit * connector.GB,
		messenger:     messenger,
		streamCfg:     messagingCfg.Stream,
		scheduler:     s,
	}, nil
}

func (s *Server) run(ctx context.Context) error {
	zap.S().Infof("Schedule reload task")
	go s.schedule()
	zap.S().Infof("Start listener ...")
	go s.listen(context.Background())
	return nil
}

// loadFromDatabase load connectors from database and run if needed
func (s *Server) loadFromDatabase() error {
	ctx := context.Background()
	zap.S().Infof("Loading connectors from db")
	connectors, err := s.connectorRepo.GetActive(ctx)
	if err != nil {
		zap.S().Errorf("Load connectors failed: %v", err)
		return err
	}
	for _, connector := range connectors {
		if err = NewTrigger(s.messenger, s.connectorRepo, s.docRepo, connector, s.fileSizeLimit).Do(ctx); err != nil {
			zap.S().Errorf("run connector %d failed: %v", connector.ID, err)
		}
	}
	return nil
}

func (s *Server) schedule() error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(s.renewInterval),
		gocron.NewTask(s.loadFromDatabase),
		gocron.WithName("reload from database"),
	)
	if err != nil {
		return err
	}
	s.scheduler.Start()
	return nil

}

// listen nats channel with updated connectors
func (s *Server) listen(ctx context.Context) {

	if err := s.loadFromDatabase(); err != nil {
		return
	}
	//if err := s.messenger.Listen(ctx, model.TopicUpdateConnector, model.SubscriptionOrchestrator, s.handleTriggerRequest); err != nil {
	//	zap.S().Errorf("failed to listen: %v", err)
	//}
}
