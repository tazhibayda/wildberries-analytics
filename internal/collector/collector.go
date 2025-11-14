package collector

import (
	"context"
	"github.com/rs/zerolog"
	"sync"
	"time"
	"wildberriesapi/internal/api"
	"wildberriesapi/internal/config"
	"wildberriesapi/internal/publisher"
)

type Collector struct {
	API          *api.WBClient
	Publisher    publisher.Publisher
	Logger       zerolog.Logger
	PollInterval time.Duration
}

func NewCollector(cfg config.Config, API *api.WBClient, pub publisher.Publisher, Logger zerolog.Logger) *Collector {
	return &Collector{
		API:          API,
		Publisher:    pub,
		Logger:       Logger,
		PollInterval: cfg.PollInterval,
	}
}

// Schedule — основной цикл периодического запуска всех сборов.
func (c *Collector) Schedule(ctx context.Context) {
	ticker := time.NewTicker(c.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Logger.Info().Msg("🚀 Starting WB full data collection cycle...")

			// Запускаем всё параллельно
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				c.CollectOrders(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				c.CollectSales(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				c.CollectStocks(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				c.collectAndPublish(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				c.collectAndPublishTarrifs(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				c.CollectAll(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				c.CollectDailyReports(ctx)
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				payload := map[string]interface{}{
					"dateFrom": time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
					"dateTo":   time.Now().Format("2006-01-02"),
				}
				c.CollectAndPublishSearchText(ctx, payload)
			}()

			wg.Wait()
			c.Logger.Info().Msg("✅ WB data collection cycle completed")

		case <-ctx.Done():
			c.Logger.Warn().Msg("🛑 Collector stopped (context cancelled)")
			return
		}
	}
}
