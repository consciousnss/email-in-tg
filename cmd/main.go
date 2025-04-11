package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/un1uckyyy/email-in-tg/internal/repo"

	"github.com/un1uckyyy/email-in-tg/internal/services/pool"
	"github.com/un1uckyyy/email-in-tg/internal/services/tg"

	"github.com/un1uckyyy/email-in-tg/pkg/slogger"

	"github.com/un1uckyyy/email-in-tg/pkg/mongo"

	"github.com/un1uckyyy/email-in-tg/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slogger.PkgLogger("main")

	cfg, err := config.LoadConfig()
	if err != nil {
		msg := fmt.Sprintf("failed to load config: %v", err)
		logger.Error(msg)
		return
	}

	db, err := mongo.New(ctx, cfg.MongoURI)
	if err != nil {
		msg := fmt.Sprintf("failed to init mongo: %v", err)
		logger.Error(msg)
		return
	}

	p := pool.NewPool()

	r := repo.NewRepo(db)
	ts, err := tg.NewTelegramService(cfg.TelegramToken, p, r)
	if err != nil {
		msg := fmt.Sprintf("failed to init telegram: %v", err)
		logger.Error(msg)
		return
	}

	err = ts.Start(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to start telegram: %v", err)
		logger.Error(msg)
		return
	}

	logger.Info("app started...")

	<-ctx.Done()
	stop()

	logger.Info("stopping gracefully...")
	ts.Stop()
}
