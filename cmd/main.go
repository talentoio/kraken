package main

import (
	"context"
	"kraken/internal/config"
	"kraken/internal/gateway"
	"kraken/internal/jobs/ltp"
	"kraken/internal/resthttp"
	"kraken/internal/services"
	"kraken/pkg/server"
	"kraken/pkg/storage/memory"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const (
	appGracefulTimeout = 30 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, _ := zap.NewProduction()

	logger.Info("service starting...")

	logger.Info("get config...")
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		logger.Error("get config", zap.Error(err))
		return
	}

	ltpMemoryCache := memory.NewCache()
	krakenGateway := gateway.NewKrakenGateway(cfg.Kraken)
	ltpService := services.NewLTPService(*cfg, logger, ltpMemoryCache, krakenGateway)
	ltpJob := ltp.NewLTPUpdaterJob(*cfg, logger, ltpMemoryCache, krakenGateway)

	ltpJob.Run(ctx)

	routeDependencies := &resthttp.RouterDependencies{
		LTPService: ltpService,
	}

	routes := resthttp.RegisterRoutes(routeDependencies)

	logger.Info("service starting...", zap.Int("port", cfg.Server.Port))
	srv, err := server.NewEchoServer(ctx, cfg.Server.Port, routes, cancel)
	if err != nil {
		logger.Error("initialize server", zap.Error(err))
		return
	}

	srv.Start()

	<-ctx.Done()
	logger.Info("signal received, stop application...")

	stopCtx, stopCancel := context.WithTimeout(context.Background(), appGracefulTimeout)
	defer stopCancel()

	// stop http server
	srv.Stop(stopCtx)

}
