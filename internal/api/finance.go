package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// FinanceOperation — структура одной финансовой операции
type FinanceOperation map[string]any

// ReturnItem — возврат товара
type ReturnItem map[string]any

// SupplyItem — поставка
type SupplyItem map[string]any

// GetFinanceOperations — аналог get_finance_operations()
func (c *WBClient) GetFinanceOperations(ctx context.Context, dateFrom, dateTo string) ([]FinanceOperation, error) {
	c.Logger.Info().Msgf("💰 Fetching finance operations %s..%s", dateFrom, dateTo)

	all := []FinanceOperation{}
	params := map[string]string{
		"dateFrom": dateFrom,
		"dateTo":   dateTo,
		"limit":    "1000",
	}

	for _, token := range c.Tokens {
		url := WBBaseURLs["finance"] + "/api/v1/supplier/finances/operations"

		body, err := c.doRequest(ctx, "GET", url, token, params)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("finance ops failed for supplier=%d")
			continue
		}

		var resp []FinanceOperation
		if err := json.Unmarshal(body, &resp); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal finance ops error")
			continue
		}

		c.Logger.Info().Msgf("✅ got %d finance operations", len(resp))
		all = append(all, resp...)
		time.Sleep(500 * time.Millisecond)
	}

	return all, nil
}

// GetReturns — аналог get_returns()
func (c *WBClient) GetReturns(ctx context.Context, dateFrom, dateTo string) ([]ReturnItem, error) {
	c.Logger.Info().Msgf("🔄 Fetching returns %s..%s", dateFrom, dateTo)

	all := []ReturnItem{}
	params := map[string]string{
		"dateFrom": dateFrom,
		"dateTo":   dateTo,
	}

	for _, token := range c.Tokens {
		url := WBBaseURLs["returns"] + "/api/v1/supplier/returns"

		body, err := c.doRequest(ctx, "GET", url, token, params)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("returns failed for supplier=%d")
			continue
		}

		var resp []ReturnItem
		if err := json.Unmarshal(body, &resp); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal returns error")
			continue
		}

		c.Logger.Info().Msgf("✅ got %d returns", len(resp))
		all = append(all, resp...)
		time.Sleep(500 * time.Millisecond)
	}

	return all, nil
}

// GetSupplies — аналог get_supplies()
func (c *WBClient) GetSupplies(ctx context.Context, limit int) ([]SupplyItem, error) {
	c.Logger.Info().Msgf("📦 Fetching supplies (limit=%d)", limit)

	all := []SupplyItem{}
	params := map[string]string{"limit": fmt.Sprint(limit)}

	for _, token := range c.Tokens {
		url := WBBaseURLs["supplies"] + "/api/v1/supplier/supplies"

		body, err := c.doRequest(ctx, "GET", url, token, params)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("supplies failed for supplier=%d")
			continue
		}

		var resp []SupplyItem
		if err := json.Unmarshal(body, &resp); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal supplies error")
			continue
		}

		c.Logger.Info().Msgf("✅  got %d supplies", len(resp))
		all = append(all, resp...)
		time.Sleep(500 * time.Millisecond)
	}

	return all, nil
}
