package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type WBRecord map[string]any

// GetOrders — получение заказов
func (c *WBClient) GetOrders(ctx context.Context, dateFrom, dateTo string) ([]WBRecord, error) {
	all := []WBRecord{}
	urlTemplate := fmt.Sprintf("https://statistics-api.wildberries.ru/api/v1/supplier/orders?dateFrom=%s", dateFrom)
	if dateTo != "" {
		urlTemplate += "&dateTo=" + dateTo
	}

	// := c.SupplierIDs[i]
	url := urlTemplate
	c.Logger.Info().Msgf("📦 Fetching orders for  from=%s to=%s", dateFrom, dateTo)

	body, err := c.doRequest(ctx, "GET", url, c.Tokens[0], nil)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("Failed to get orders for ")
		return nil, err
	}

	var data []WBRecord
	if err := json.Unmarshal(body, &data); err != nil {
		c.Logger.Error().Err(err).Msg("unmarshal orders error")
		return nil, err
	}

	all = append(all, data...)

	return all, nil
}

// GetSales — получение продаж
func (c *WBClient) GetSales(ctx context.Context, dateFrom string) ([]WBRecord, error) {
	all := []WBRecord{}
	urlTemplate := fmt.Sprintf("https://statistics-api.wildberries.ru/api/v1/supplier/sales?dateFrom=%s", dateFrom)

	for _, token := range c.Tokens {
		c.Logger.Info().Msgf("💰 Fetching sales for  from=%s", dateFrom)

		body, err := c.doRequest(ctx, "GET", urlTemplate, token, nil)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("Failed to get sales for ")
			continue
		}

		var data []WBRecord
		if err := json.Unmarshal(body, &data); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal sales error")
			continue
		}

		all = append(all, data...)
	}

	return all, nil
}

// GetStocks — получение остатков
func (c *WBClient) GetStocks(ctx context.Context, dateFrom string) ([]WBRecord, error) {
	all := []WBRecord{}
	urlTemplate := fmt.Sprintf("https://statistics-api.wildberries.ru/api/v1/supplier/stocks?dateFrom=%s", dateFrom)

	for _, token := range c.Tokens {
		c.Logger.Info().Msgf("📦 Fetching stocks for  from=%s", dateFrom)

		body, err := c.doRequest(ctx, "GET", urlTemplate, token, nil)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("Failed to get stocks for ")
			continue
		}

		var data []WBRecord
		if err := json.Unmarshal(body, &data); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal stocks error")
			continue
		}

		all = append(all, data...)
	}

	return all, nil
}

// GetIncomes — получение остатков
func (c *WBClient) GetIncomes(ctx context.Context, dateFrom string) ([]WBRecord, error) {
	all := []WBRecord{}
	urlTemplate := fmt.Sprintf("https://statistics-api.wildberries.ru/api/v1/supplier/incomes?dateFrom=%s", dateFrom)

	for _, token := range c.Tokens {
		c.Logger.Info().Msgf("📦 Fetching incomes for  from=%s", dateFrom)

		body, err := c.doRequest(ctx, "GET", urlTemplate, token, nil)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("Failed to get incomes for ")
			continue
		}

		var data []WBRecord
		if err := json.Unmarshal(body, &data); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal incomes error")
			continue
		}

		all = append(all, data...)
	}

	return all, nil
}
