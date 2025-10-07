package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Sale struct {
	NmID            int       `json:"nm_id"`
	SupplierArticle string    `json:"supplier_article"`
	Quantity        int       `json:"quantity"`
	TotalPrice      float64   `json:"total_price"`
	Date            time.Time `json:"date"`
	LastChangeDate  time.Time `json:"last_change_date"`
}

func (c *Collector) CollectSales() error {
	ctx := context.Background()

	sales, err := c.API.GetSales(ctx, time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get sales")
		return err
	}

	for _, s := range sales {
		data, _ := json.Marshal(s)
		if err := c.Publisher.Publish(ctx, "wb.raw.sales", nil, data); err != nil {
			c.Logger.Error().Err(err).Msg("failed to publish sale: %v")
		}
		fmt.Println(data)
	}
	c.Logger.Info().Msg(fmt.Sprintf("✅ published %d sales", len(sales)))
	return nil
}
