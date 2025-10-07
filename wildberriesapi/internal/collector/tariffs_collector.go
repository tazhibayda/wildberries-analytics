package collector

import (
	"context"
	"encoding/json"
	"fmt"
)

type Tariff struct {
	RegionName   string  `json:"region_name"`
	Price        float64 `json:"price"`
	DeliveryDays int     `json:"delivery_days"`
}

func (c *Collector) CollectTariffs() error {
	ctx := context.Background()
	tariffs, err := c.API.GetTariffs(ctx)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get tariffs: %v")
		return err
	}

	for _, t := range tariffs {
		data, _ := json.Marshal(t)
		if err := c.Publisher.Publish(ctx, "wb.raw.tariffs", nil, data); err != nil {
			c.Logger.Error().Err(err).Msg("failed to publish tariff: %v")
		}
	}
	c.Logger.Info().Msg(fmt.Sprintf("✅ published %d tariffs", len(tariffs)))
	return nil
}
