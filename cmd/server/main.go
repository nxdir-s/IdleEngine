package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nxdir-s/IdleEngine/internal/adapters/secondary"
	"github.com/nxdir-s/IdleEngine/internal/config"
	"github.com/nxdir-s/IdleEngine/internal/engine"
	"github.com/nxdir-s/IdleEngine/internal/ports"
	"github.com/nxdir-s/IdleEngine/internal/server"
	"go.opentelemetry.io/otel"
	"golang.org/x/sys/unix"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		logger.Error("failed to get RLIMIT", slog.Any("err", err))
		os.Exit(1)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		logger.Error("failed to set RLIMIT", slog.Any("err", err))
		os.Exit(1)
	}

	cfg, err := config.New(
		config.WithListenerAddr(),
		config.WithBrokers(),
		config.WithUserEventsTopic(),
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// otelCfg := &telemetry.Config{
	// 	ServiceName:  cfg.OtelService,
	// 	OtelEndpoint: cfg.OtelEndpoint,
	// 	Insecure:     true,
	// }

	// cleanup, err := telemetry.InitProviders(ctx, otelCfg)
	// if err != nil {
	// 	logger.Error("failed to initialize telemetry", slog.Any("err", err))
	// 	os.Exit(1)
	// }
	// defer cleanup(ctx)

	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", cfg.ListenerAddr)
	if err != nil {
		logger.Error("failed to create tcp listener", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("listening on ws://%v", listener.Addr()))

	var kafka ports.Kafka
	kafka, err = secondary.NewFranzAdapter(
		logger,
		otel.Tracer("kafka.franz"),
		secondary.WithProducer(cfg.UserEventsTopic, strings.Split(cfg.Brokers, ",")),
	)
	if err != nil {
		logger.Error("failed to create kafka adapter", slog.Any("err", err))
		os.Exit(1)
	}
	defer kafka.Close()

	fd, err := unix.EpollCreate1(0)
	if err != nil {
		logger.Error("failed to create epoll", slog.Any("err", err))
		os.Exit(1)
	}

	pool := server.NewPool(ctx, logger, otel.Tracer("pool"))
	epoll := server.NewEpoll(fd, pool, logger, otel.Tracer("epoll"))
	ngin := engine.NewGameEngine(pool, kafka, logger, otel.Tracer("engine"))

	server := server.NewGameServer(ctx, listener, epoll, ngin, pool, logger)

	logger.Info("starting server")

	go server.Start(ctx)

	select {
	case <-ctx.Done():
		logger.Warn(ctx.Err().Error())
	}

	if err := server.Shutdown(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
