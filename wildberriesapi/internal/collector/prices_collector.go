package collector

import (
	"context"
	"encoding/json"
	"fmt"
)

type Price struct {
	NmID            int     `json:"nm_id"`
	SupplierArticle string  `json:"supplier_article"`
	Price           float64 `json:"price"`
	DiscountPercent int     `json:"discount_percent"`
	Currency        string  `json:"currency"`
}

func (c *Collector) CollectPrices() error {
	ctx := context.Background()
	prices, err := c.API.GetPrices(ctx)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get prices: %v")
		return err
	}

	for _, p := range prices {
		data, _ := json.Marshal(p)
		if err := c.Publisher.Publish(ctx, "wb.raw.prices", nil, data); err != nil {
			c.Logger.Error().Err(err).Msg("failed to publish price: %v")
		}
	}
	c.Logger.Info().Msg(fmt.Sprintf("✅ published %d tariffs", len(prices)))
	return nil
}
