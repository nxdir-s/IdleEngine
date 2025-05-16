package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	DefaultMaxClients int = 50
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if len(os.Args) < 2 {
		logger.Error("please provide an address")
		os.Exit(1)
	}

	var maxClients int
	maxClients = DefaultMaxClients

	if len(os.Args) == 3 {
		maxClientArg, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		maxClients = maxClientArg
	}

	var wg sync.WaitGroup
	for range maxClients {
		time.Sleep(80 * time.Millisecond)

		wg.Add(1)
		go startWS(ctx, &wg, os.Args[1], logger)
	}

	wg.Wait()
}

func startWS(ctx context.Context, wg *sync.WaitGroup, addr string, logger *slog.Logger) {
	defer wg.Done()

	conn, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		logger.Error("failed to dial websocket", slog.String("err", err.Error()))
		return
	}
	defer conn.CloseNow()

	for {
		select {
		case <-ctx.Done():
			conn.Close(websocket.StatusNormalClosure, "")
			return
		default:
			msgType, r, err := conn.Reader(ctx)
			if err != nil {
				logger.Error("failed to get reader from connection", slog.String("err", err.Error()))
				return
			}

			var msg []byte
			msg, err = io.ReadAll(r)
			if err != nil {
				logger.Error("failed reading message", slog.String("err", err.Error()))
				return
			}

			switch msgType {
			case websocket.MessageText:
				logger.Info(string(msg))
			case websocket.MessageBinary:
				logger.Info("recieved binary message")
			default:
				logger.Warn("unknown websocket frame type", slog.Any("msgType", msgType))
				return
			}
		}
	}
}
