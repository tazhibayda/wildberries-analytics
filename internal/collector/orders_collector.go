package collector

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"wildberriesapi/internal/api"
	"wildberriesapi/internal/publisher"
)

type OrdersCollector struct {
	api     *api.WBClient
	pub     publisher.Publisher
	Logger  zerolog.Logger
	topic   string
	enabled bool
}

func (c *Collector) CollectOrders(ctx context.Context) {

	dateFrom := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	dateTo := time.Now().Format("2006-01-02")

	c.Logger.Info().Msgf("📦 Collecting WB orders from %s to %s", dateFrom, dateTo)
	data, err := c.API.GetOrders(ctx, dateFrom, dateTo)
	if err != nil {
		c.Logger.Error().Err(err).Msg("❌ Failed to collect WB orders")
		return
	}

	if len(data) == 0 {
		c.Logger.Info().Msg("📦 No new WB orders found")
		return
	}

	if err := c.Publisher.Publish(ctx, "wb.raw.orders", []byte("orders"), data); err != nil {
		c.Logger.Error().Err(err).Msg("❌ Failed to publish WB orders to Kafka")
		return
	}

	c.Logger.Info().Msgf("✅ Published %d WB orders to topic '%s'", len(data), "wb.raw.orders")
}
