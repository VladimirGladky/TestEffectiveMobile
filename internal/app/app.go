package app

import (
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/internal/repository"
	"TestEffectiveMobile/internal/service"
	"TestEffectiveMobile/internal/transport"
	"TestEffectiveMobile/pkg/logger"
	"TestEffectiveMobile/pkg/postgres"
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	SubscriptionServer *transport.SubscriptionServer
	cfg                *config.Config
	ctx                context.Context
	wg                 sync.WaitGroup
	cancel             context.CancelFunc
}

func New(cfg *config.Config, ctx context.Context) *App {
	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	repo := repository.NewSubscriptionRepository(db, ctx)
	srv := service.NewSubscriptionService(repo, cfg, ctx)
	server := transport.New(srv, cfg, ctx)
	return &App{
		SubscriptionServer: server,
		cfg:                cfg,
		ctx:                ctx,
	}
}

func (s *App) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	a.wg.Add(1)
	go func() {
		logger.GetLoggerFromCtx(a.ctx).Info("Server started on address", zap.Any("address", a.cfg.Host+":"+a.cfg.Port))
		defer a.wg.Done()
		if err := a.SubscriptionServer.Run(); err != nil {
			errCh <- err
			a.cancel()
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errCh:
		logger.GetLoggerFromCtx(a.ctx).Error("error running app", zap.Error(err))
		return err
	case <-a.ctx.Done():
		logger.GetLoggerFromCtx(a.ctx).Info("context done")
	}

	return nil
}
