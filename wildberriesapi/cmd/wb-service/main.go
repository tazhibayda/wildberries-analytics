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
	// --- 1️⃣ Загрузка конфигурации ---
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	log.Info().Msg("🚀 Starting WB Analytics Collector Service")

	// --- 2️⃣ Создаём общий контекст с отменой ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- 3️⃣ Инициализация клиентов ---
	wbClient := api.NewWBClient(cfg)

	pub, err := publisher.NewKafkaPublisher(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Failed to create Kafka publisher")
	}
	defer pub.Close()

	// --- 4️⃣ Создаём коллектор ---
	coll := collector.NewCollector(cfg, wbClient, pub, log)

	// --- 5️⃣ Запуск планировщика ---
	go func() {
		log.Info().Msgf("⏱️ Collector scheduler started (interval: %s)", cfg.PollInterval)
		coll.Schedule(ctx, cfg.PollInterval)
	}()

	// --- 6️⃣ Graceful Shutdown ---
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		log.Warn().Msg("🛑 Shutdown signal received, stopping service...")
		cancel()
		time.Sleep(2 * time.Second)
	}

	log.Info().Msg("✅ Service stopped gracefully")
}
