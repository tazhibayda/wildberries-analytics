package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Stock struct {
	NmID            int       `json:"nm_id"`
	SupplierArticle string    `json:"supplier_article"`
	WarehouseName   string    `json:"warehouse_name"`
	Quantity        int       `json:"quantity"`
	LastChangeDate  time.Time `json:"last_change_date"`
}

func (c *Collector) CollectStocks() error {
	ctx := context.Background()
	stocks, err := c.API.GetStocks(ctx, time.Now().Add(-24*time.Hour))
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get stocks")
		return err
	}

	for _, s := range stocks {
		data, _ := json.Marshal(s)
		if err := c.Publisher.Publish(ctx, "wb.raw.stocks", nil, data); err != nil {
			c.Logger.Error().Err(err).Msg("failed to publish stock")
		}
	}
	c.Logger.Info().Msg(fmt.Sprintf("✅ published %d stocks", len(stocks)))
	return nil
}
