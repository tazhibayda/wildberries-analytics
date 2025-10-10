package collector

import (
	"context"
)

func (c *Collector) collectAndPublish(ctx context.Context) {
	c.Logger.Info().Msg("📥 Collecting prices from WB API")

	prices, err := c.API.GetPrices(ctx, 1000, 0)
	if err != nil {
		c.Logger.Error().Err(err).Msg("❌ failed to fetch prices")
		return
	}

	count := 0
	for _, p := range prices {
		err := c.Publisher.Publish(ctx, ("wb.raw.prices"), nil, p)
		if err == nil {
			count++
		}
	}

	c.Logger.Info().Msgf("✅ Published %d price records to Kafka topic 'wb.prices'", count)
}
