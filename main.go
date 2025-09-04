package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	mongoinfra "github.com/un1uckyyy/email-in-tg/internal/infra/mongo"

	"github.com/un1uckyyy/email-in-tg/internal/app/pool"
	"github.com/un1uckyyy/email-in-tg/internal/app/tg"

	"github.com/un1uckyyy/email-in-tg/pkg/otel"
	"github.com/un1uckyyy/email-in-tg/pkg/slogger"

	"github.com/un1uckyyy/email-in-tg/pkg/mongo"

	"github.com/un1uckyyy/email-in-tg/internal/config"
)

const appName = "email-in-tg"

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slogger.PkgLogger("main")

	cfg, err := config.LoadConfig()
	if err != nil {
		msg := fmt.Sprintf("failed to load config: %v", err)
		logger.Error(msg)
		return err
	}

	if cfg.OtelURI != "" {
		otelShutdown, e := otel.SetupOTelSDK(ctx, appName, cfg.OtelURI)
		if e != nil {
			return e
		}

		defer func() {
			err = errors.Join(err, otelShutdown(context.Background()))
		}()
	}

	db, err := mongo.New(ctx, cfg.MongoURI)
	if err != nil {
		msg := fmt.Sprintf("failed to init mongo: %v", err)
		logger.Error(msg)
		return err
	}

	p := pool.NewPool()

	groupRepo := mongoinfra.NewGroupRepo(db.DB)
	subRepo := mongoinfra.NewSubscriptionRepo(db.DB)

	ts, err := tg.NewTelegramService(cfg.TelegramToken, p, groupRepo, subRepo)
	if err != nil {
		msg := fmt.Sprintf("failed to init telegram: %v", err)
		logger.Error(msg)
		return err
	}

	err = ts.Start(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to start telegram: %v", err)
		logger.Error(msg)
		return err
	}

	logger.Info("app started...")

	<-ctx.Done()
	stop()

	logger.Info("stopping gracefully...")
	ts.Stop()
	return nil
}
