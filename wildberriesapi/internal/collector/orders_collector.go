package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Order struct {
	OrderID         int       `json:"order_id"`
	NmID            int       `json:"nm_id"`
	SupplierArticle string    `json:"supplier_article"`
	TechSize        string    `json:"tech_size"`
	WarehouseName   string    `json:"warehouse_name"`
	TotalPrice      float64   `json:"total_price"`
	DiscountPercent int       `json:"discount_percent"`
	Date            time.Time `json:"date"`
	LastChangeDate  time.Time `json:"last_change_date"`
}

func (c *Collector) CollectOrders() error {
	ctx := context.Background()
	orders, err := c.API.GetOrders(ctx, (time.Now().Add(-1 * time.Hour)), time.Now())
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to get orders")
		return err
	}

	for _, o := range orders {
		data, _ := json.Marshal(o)
		if err := c.Publisher.Publish(ctx, "wb.raw.orders", nil, data); err != nil {
			c.Logger.Error().Err(err).Msg("failed to publish order")
		}
	}
	c.Logger.Info().Msg(fmt.Sprintf("✅ published %d tariffs", len(orders)))
	return nil
}
