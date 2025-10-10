package api

import (
	"context"
	"encoding/json"
	"time"
)

// AdvertCampaign — структура для описания рекламной кампании WB
type AdvertCampaign struct {
	CampaignID   int64   `json:"advertId"`
	CampaignName string  `json:"name"`
	Status       string  `json:"status"`
	Budget       float64 `json:"dailyBudget"`
	Type         string  `json:"type"`
	SupplierID   int     `json:"__supplier_id"`
}

// GetAdverts получает список всех рекламных кампаний по токенам
func (c *WBClient) GetAdverts(ctx context.Context) ([]AdvertCampaign, error) {
	allCampaigns := make([]AdvertCampaign, 0)

	for _, token := range c.Tokens {
		if token == "" {
			continue
		}

		url := "https://advert-api.wildberries.ru/adv/v1/promotion/count"

		body, err := c.doRequest(ctx, "GET", url, token, nil)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("❌ failed to fetch adverts ")
			continue
		}

		var resp struct {
			Data []AdvertCampaign `json:"adverts"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal error in GetAdverts response")
			continue
		}

		allCampaigns = append(allCampaigns, resp.Data...)
		c.Logger.Info().Msgf("✅ supplier_id=%d: adverts loaded (%d records)", len(resp.Data))

		time.Sleep(700 * time.Millisecond)
	}

	return allCampaigns, nil
}
