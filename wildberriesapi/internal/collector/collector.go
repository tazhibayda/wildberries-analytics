package collector

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"wildberriesapi/internal/api"
	"wildberriesapi/internal/config"
	"wildberriesapi/internal/models"
	"wildberriesapi/internal/publisher"
)

type Collector struct {
	cfg    config.Config
	client *api.WBClient
	pub    publisher.Publisher
	log    zerolog.Logger
}

func NewCollector(cfg config.Config, client *api.WBClient, pub publisher.Publisher, log zerolog.Logger) *Collector {
	return &Collector{
		cfg:    cfg,
		client: client,
		pub:    pub,
		log:    log,
	}
}

// Schedule runs periodic job
func (c *Collector) Schedule(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Первый запуск сразу
	c.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			c.log.Info().Msg("collector: context cancelled")
			return
		case <-ticker.C:
			c.runOnce(ctx)
		}
	}
}

func (c *Collector) runOnce(ctx context.Context) {
	c.log.Info().Msg("collector: run once")
	dateFrom := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

	// sales
	if sales, err := c.client.GetSales(ctx, dateFrom); err != nil {
		c.log.Error().Err(err).Msg("get sales")
	} else {
		c.publishEvents(ctx, "sales", sales)
	}

	// stocks
	if stocks, err := c.client.GetStocks(ctx, dateFrom); err != nil {
		c.log.Error().Err(err).Msg("get stocks")
	} else {
		c.publishEvents(ctx, "stocks", stocks)
	}

	// orders
	if orders, err := c.client.GetOrders(ctx, dateFrom, ""); err != nil {
		c.log.Error().Err(err).Msg("get orders")
	} else {
		c.publishEvents(ctx, "orders", orders)
	}
}

func (c *Collector) publishEvents(ctx context.Context, eventType string, items []map[string]interface{}) {
	for _, s := range items {
		ev := models.WBEvent{
			Type:       eventType,
			SupplierID: intValueOrZero(s["__supplier_id"]),
			Data:       s,
			CreatedAt:  time.Now().Format(time.RFC3339),
			Source:     "wildberries",
		}
		if err := c.pub.Publish(ctx, nil, ev); err != nil {
			c.log.Error().Err(err).Msgf("failed to publish %s event", eventType)
		}
	}
	c.log.Info().Msgf("published %d %s events", len(items), eventType)
}

func intValueOrZero(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		// можно добавить парсинг строки в число, если нужно
		return 0
	default:
		return 0
	}
}
