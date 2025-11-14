package collector

import (
	"context"
	"time"
)

func (c *Collector) CollectSales(ctx context.Context) {

	dateFrom := time.Now().Add(-24 * time.Hour).Format("2006-01-02T00:00:00")

	c.Logger.Info().Msgf("💰 Collecting WB sales since %s", dateFrom)
	data, err := c.API.GetSales(ctx, dateFrom)
	if err != nil {
		c.Logger.Error().Err(err).Msg("❌ Failed to collect WB sales")
		return
	}

	if len(data) == 0 {
		c.Logger.Info().Msg("💰 No new WB sales found")
		return
	}

	if err := c.Publisher.Publish(ctx, "wb.raw.sales", []byte("sales"), data); err != nil {
		c.Logger.Error().Err(err).Msg("❌ Failed to publish WB sales to Kafka")
		return
	}

	c.Logger.Info().Msgf("✅ Published %d WB sales to topic '%s'", len(data), "wb.raw.sales")
}
