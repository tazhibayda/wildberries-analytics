package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wildberriesapi/internal/api"
	"wildberriesapi/internal/collector"
	"wildberriesapi/internal/config"
	"wildberriesapi/internal/logger"
	"wildberriesapi/internal/publisher"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)
	log.Info().Msg("starting wb-data-parser service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create components
	wbClient := api.NewWBClient(cfg)

	pub, err := publisher.NewKafkaPublisher(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create kafka publisher")
	}

	collector := collector.NewCollector(cfg, wbClient, pub, log)

	go collector.Schedule(ctx, cfg.PollInterval)

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info().Msg("shutdown signal received")
	cancel()
	time.Sleep(2 * time.Second)
	_ = pub.Close()
	log.Info().Msg("service stopped")
}
