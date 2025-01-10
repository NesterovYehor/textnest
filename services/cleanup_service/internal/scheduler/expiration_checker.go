package scheduler

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/services"
)

type Checker struct {
	service *services.ExpirationService
	log     *jsonlog.Logger
}

func NewChecker(service *services.ExpirationService, log *jsonlog.Logger) *Checker {
	return &Checker{
		service: service,
		log:     log,
	}
}

func (s *Checker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	s.log.PrintInfo(ctx, "Scheduler started", nil)

	for {
		select {
		case <-ticker.C:
			if err := s.service.ProcessExpirations(ctx); err != nil {
				s.log.PrintError(ctx, err, nil)
			}

		case <-stopSignal:
			s.log.PrintInfo(ctx, "Received shutdown signal, stopping scheduler", nil)
			return
		}
	}
}
