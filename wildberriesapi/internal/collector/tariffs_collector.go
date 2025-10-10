package collector

import (
	"context"
)

func (c *Collector) collectAndPublishTarrifs(ctx context.Context) {
	c.Logger.Info().Msg("📥 Collecting tariffs from WB API")

	tariffs, err := c.API.GetTariffs(ctx)
	if err != nil {
		c.Logger.Error().Err(err).Msg("❌ failed to fetch tariffs")
		return
	}

	count := 0
	for _, t := range tariffs {
		err := c.Publisher.Publish(ctx, "wb.raw.tariffs", nil, t)
		if err == nil {
			count++
		}
	}

	c.Logger.Info().Msgf("✅ Published %d tariff records to Kafka topic 'wb.tariffs'", count)
}
