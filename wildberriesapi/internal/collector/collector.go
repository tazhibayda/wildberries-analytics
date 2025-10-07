package collector

import (
	"context"
	"github.com/rs/zerolog"
	"time"
	"wildberriesapi/internal/api"
	"wildberriesapi/internal/config"
	"wildberriesapi/internal/publisher"
)

type Collector struct {
	API       *api.WBClient
	Publisher publisher.Publisher
	Logger    zerolog.Logger
}

func NewCollector(cfg config.Config, API *api.WBClient, pub publisher.Publisher, Logger zerolog.Logger) *Collector {
	return &Collector{
		API:       API,
		Publisher: pub,
		Logger:    Logger,
	}
}

func (c *Collector) Schedule(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Logger.Info().Msg("🚀 Running WB data collection cycle...")
			c.CollectOrders()
			c.CollectSales()
			c.CollectStocks()
			c.CollectPrices()
			c.CollectTariffs()
			c.Logger.Info().Msg("✅ WB data collection cycle completed")
		case <-ctx.Done():
			c.Logger.Warn().Msg("🛑 Collector stopped by context cancel")
			return
		}
	}
}
