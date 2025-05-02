package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/nxdir-s/IdleEngine/internal/adapters/primary"
	"github.com/nxdir-s/IdleEngine/internal/adapters/secondary"
	"github.com/nxdir-s/IdleEngine/internal/config"
	"github.com/nxdir-s/IdleEngine/internal/consumers"
	"github.com/nxdir-s/IdleEngine/internal/core/domain"
	"github.com/nxdir-s/IdleEngine/internal/core/service"
	"github.com/nxdir-s/IdleEngine/internal/logs"
	"github.com/nxdir-s/IdleEngine/internal/observability"
	"github.com/nxdir-s/IdleEngine/internal/ports"
	"github.com/nxdir-s/telemetry"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(logs.NewHandler(slog.NewTextHandler(os.Stdout, nil)))
	slog.SetDefault(logger)

	cfg, err := config.New(
		config.WithBrokers(),
		config.WithRedPandaUsr(),
		config.WithRedPandaPass(),
		config.WithConsumerName(),
		config.WithOtelServiceName(),
		config.WithOtelEndpoint(),
		config.WithProfileURL(),
		config.WithGrafanaUsr(),
		config.WithGrafanaPass(),
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	otelCfg := &telemetry.Config{
		ServiceName:        cfg.OtelService,
		OtelEndpoint:       cfg.OtelEndpoint,
		Insecure:           true,
		EnableSpanProfiles: true,
	}

	cleanup, err := telemetry.InitProviders(ctx, otelCfg)
	if err != nil {
		logger.Error("failed to initialize telemetry", slog.Any("err", err))
		os.Exit(1)
	}
	defer cleanup(ctx)

	profileCfg := &observability.ProfileConfig{
		ApplicationName: cfg.OtelService,
		ServerAddress:   cfg.ProfileURL,
		AuthUser:        cfg.GrafanaUsr,
		AuthPassword:    cfg.GrafanaPass,
	}

	profiler, err := observability.NewProfiler(profileCfg)
	if err != nil {
		logger.Error("failed to start profiler", slog.Any("err", err))
		os.Exit(1)
	}
	defer profiler.Stop()

	var pgxPool secondary.PgxPool
	pgxPool, err = secondary.NewPgxPool(ctx, "dbUrl")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// secondary adapters
	var database ports.Database
	var kafka ports.Kafka

	// services
	var userService ports.UserService

	// domain orchestrators
	var events ports.Events
	var users ports.Users

	// primary adapter
	var adapter ports.Consumer

	kafka, err = secondary.NewFranzAdapter(
		logger,
		otel.Tracer("kafka.franz"),
		secondary.WithConsumer(
			"user.events",
			cfg.ConsumerName,
			strings.Split(cfg.Brokers, ","),
			cfg.RedPandaUsr,
			cfg.RedPandaPass,
		),
	)
	if err != nil {
		logger.Error("failed to create kafka adapter", slog.Any("err", err))
		os.Exit(1)
	}

	database = secondary.NewPostgresAdapter(pgxPool, logger, otel.Tracer("postgres"))

	userService = service.NewUserService(database)

	users = domain.NewUsers(userService)
	events = domain.NewEvents(users)

	adapter = primary.NewConsumerAdapter(events, otel.Tracer("consumer.userevents"))

	consumer := consumers.NewUserEvents(kafka, adapter, logger)
	defer consumer.Close()

	go consumer.Start(ctx)

	select {
	case <-ctx.Done():
		logger.Info(ctx.Err().Error())
		os.Exit(0)
	}
}
