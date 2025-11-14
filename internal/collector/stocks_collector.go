package collector

import (
	"context"
	"time"
)

func (c *Collector) CollectStocks(ctx context.Context) {
	dateFrom := time.Now().Add(-24 * time.Hour).Format("2006-01-02T00:00:00")

	c.Logger.Info().Msgf("📊 Collecting WB stocks since %s", dateFrom)
	data, err := c.API.GetStocks(ctx, dateFrom)
	if err != nil {
		c.Logger.Error().Err(err).Msg("❌ Failed to collect WB stocks")
		return
	}

	if len(data) == 0 {
		c.Logger.Info().Msg("📊 No new WB stocks found")
		return
	}

	if err := c.Publisher.Publish(ctx, "wb.raw.stocks", []byte("stocks"), data); err != nil {
		c.Logger.Error().Err(err).Msg("❌ Failed to publish WB stocks to Kafka")
		return
	}

	c.Logger.Info().Msgf("✅ Published %d WB stocks to topic '%s'", len(data), "wb.raw.stocks")
}
