package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"wildberriesapi/internal/config"
)

// WBClient - minimal WB API client
type WBClient struct {
	token   string
	timeout time.Duration
	client  *http.Client
	base    string
}

func NewWBClient(cfg config.Config) *WBClient {
	return &WBClient{
		token:   cfg.WBToken,
		timeout: cfg.HTTPTimeout,
		client:  &http.Client{Timeout: cfg.HTTPTimeout},
		base:    "https://statistics-api.wildberries.ru/api/v1/supplier",
	}
}

func (c *WBClient) doGet(ctx context.Context, url string) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// parseResponse ensures we always return []map[string]interface{}
func parseResponse(raw interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	switch x := raw.(type) {
	case []interface{}:
		for _, v := range x {
			if m, ok := v.(map[string]interface{}); ok {
				result = append(result, m)
			} else {
				return nil, fmt.Errorf("unexpected element type: %T", v)
			}
		}
	case map[string]interface{}:
		result = append(result, x)
	default:
		return nil, fmt.Errorf("unexpected response type: %T", raw)
	}
	return result, nil
}

// GetSales - example: dateFrom in RFC3339 (string)
func (c *WBClient) GetSales(ctx context.Context, dateFrom string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/sales?dateFrom=%s", c.base, dateFrom)
	raw, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseResponse(raw)
}

// GetStocks
func (c *WBClient) GetStocks(ctx context.Context, dateFrom string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/stocks?dateFrom=%s", c.base, dateFrom)
	raw, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseResponse(raw)
}

// GetOrders
func (c *WBClient) GetOrders(ctx context.Context, dateFrom, dateTo string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/orders?dateFrom=%s", c.base, dateFrom)
	if dateTo != "" {
		url += "&dateTo=" + dateTo
	}
	raw, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseResponse(raw)
}
