package main

import (
	"net/http"
	_ "wildberriesapi/docs"
	"wildberriesapi/internal/api"
	"wildberriesapi/internal/config"
	"wildberriesapi/internal/handlers"
	"wildberriesapi/internal/logger"
)

// @title WB Analytics Collector Service API
// @version 1.0
// @description This is the API documentation for the WB Analytics Collector Service.
// @BasePath        /
// @schemes         http https

// @in header
func main() {
	// --- 1️⃣ Загрузка конфигурации ---
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	log.Info().Msg("🚀 Starting WB Analytics Collector Service")

	// --- 2️⃣ Создаём общий контекст с отменой ---
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	tokens := make([]string, 0)
	tokens = append(tokens, cfg.WBToken)
	// --- 3️⃣ Инициализация клиентов ---
	wbClient := api.NewWBClient(tokens, log)

	handler := handlers.NewRouter(wbClient, log)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return
	}

	//pub, err := publisher.NewKafkaPublisher(cfg)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("❌ Failed to create Kafka publisher")
	//}
	//defer pub.Close()
	//
	//// --- 4️⃣ Создаём коллектор ---
	//coll := collector.NewCollector(cfg, wbClient, pub, log)
	//
	//// --- 5️⃣ Запуск планировщика ---
	//go func() {
	//	log.Info().Msgf("⏱️ Collector scheduler started (interval: %s)", cfg.PollInterval)
	//	coll.Schedule(ctx)
	//}()
	//
	//// --- 6️⃣ Graceful Shutdown ---
	//sig := make(chan os.Signal, 1)
	//signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	//
	//select {
	//case <-sig:
	//	log.Warn().Msg("🛑 Shutdown signal received, stopping service...")
	//	cancel()
	//	time.Sleep(2 * time.Second)
	//}

	log.Info().Msg("✅ Service stopped gracefully")
}
