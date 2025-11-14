package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// PriceItem — структура одной записи о товаре из WB API
type PriceItem struct {
	ID          int64   `json:"nmId"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	SupplierID  int     `json:"__supplier_id"`
	SupplierArt string  `json:"supplierArticle,omitempty"`
	// можно добавить другие поля по необходимости
}

// GetPrices получает список товаров с ценами постранично по каждому токену
func (c *WBClient) GetPrices(ctx context.Context, limit, offset int) ([]PriceItem, error) {
	allPrices := make([]PriceItem, 0)

	for _, token := range c.Tokens {
		if token == "" {
			continue
		}
		tokenTotal := 0
		pageOffset := offset

		for {
			params := map[string]string{
				"limit":  fmt.Sprintf("%d", limit),
				"offset": fmt.Sprintf("%d", pageOffset),
			}

			body, err := c.doRequest(ctx, "GET", WBEndpoints.Prices.URL, token, params)
			if err != nil {
				c.Logger.Error().Err(err).Msgf("❌ failed to fetch prices ( offset=%d)", pageOffset)
				break
			}

			var resp struct {
				Data struct {
					ListGoods []PriceItem `json:"listGoods"`
				} `json:"data"`
			}

			if err := json.Unmarshal(body, &resp); err != nil {
				c.Logger.Error().Err(err).Msg("unmarshal error in GetPrices response")
				break
			}

			goods := resp.Data.ListGoods
			if len(goods) == 0 {
				c.Logger.Info().Msgf("ℹ️ Empty response ( offset=%d) — stopping", pageOffset)
				break
			}

			allPrices = append(allPrices, goods...)
			tokenTotal += len(goods)

			c.Logger.Info().Msgf("📦 : fetched %d goods (offset=%d)", len(goods), pageOffset)

			// мягкий rate-limit WB API
			time.Sleep(600 * time.Millisecond)

			// если вернулось меньше лимита — значит это последняя страница
			if len(goods) < limit {
				break
			}
			pageOffset += limit
		}

		c.Logger.Info().Msgf("✅ : total %d price records collected", tokenTotal)
	}

	return allPrices, nil
}
