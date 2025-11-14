package collector

import (
	"context"
	"encoding/json"
	"time"
)

func (c *Collector) CollectAll(ctx context.Context) {
	dateTo := time.Now().Format("2006-01-02")
	dateFrom := time.Now().AddDate(0, 0, -7).Format("2006-01-02") // последние 7 дней

	c.Logger.Info().Msgf("🏦 Starting finance collection %s..%s", dateFrom, dateTo)

	// --- 1️⃣ Финансовые операции
	ops, err := c.API.GetFinanceOperations(ctx, dateFrom, dateTo)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to collect finance operations")
	} else {
		payload := map[string]any{
			"type":      "finance_operations",
			"dateFrom":  dateFrom,
			"dateTo":    dateTo,
			"timestamp": time.Now().Format(time.RFC3339),
			"data":      ops,
		}
		b, _ := json.Marshal(payload)
		c.Publisher.Publish(ctx, "wb.finance.operations", nil, b)
	}

	// --- 2️⃣ Возвраты
	returns, err := c.API.GetReturns(ctx, dateFrom, dateTo)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to collect returns")
	} else {
		payload := map[string]any{
			"type":      "returns",
			"dateFrom":  dateFrom,
			"dateTo":    dateTo,
			"timestamp": time.Now().Format(time.RFC3339),
			"data":      returns,
		}
		b, _ := json.Marshal(payload)
		c.Publisher.Publish(ctx, "wb.raw.finance.returns", nil, b)
	}

	// --- 3️⃣ Поставки
	supplies, err := c.API.GetSupplies(ctx, 1000)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to collect supplies")
	} else {
		payload := map[string]any{
			"type":      "supplies",
			"timestamp": time.Now().Format(time.RFC3339),
			"data":      supplies,
		}
		b, _ := json.Marshal(payload)
		c.Publisher.Publish(ctx, "wb.raw.finance.supplies", nil, b)
	}

	c.Logger.Info().Msg("✅ Finance data collection completed")
}
