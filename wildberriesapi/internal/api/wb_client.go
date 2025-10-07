package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"wildberriesapi/internal/config"
)

// WBClient — клиент для Wildberries API
type WBClient struct {
	token   string
	timeout time.Duration
	client  *http.Client
	base    string
}

// NewWBClient создаёт новый WB API клиент
func NewWBClient(cfg config.Config) *WBClient {
	return &WBClient{
		token:   cfg.WBToken,
		timeout: cfg.HTTPTimeout,
		client:  &http.Client{Timeout: cfg.HTTPTimeout},
		base:    "https://statistics-api.wildberries.ru/api/v1/supplier",
	}
}

// retry wrapper с экспоненциальной задержкой
func retry(attempts int, sleep time.Duration, fn func() error) error {
	err := fn()
	if err == nil {
		return nil
	}

	for i := 1; i < attempts; i++ {
		time.Sleep(sleep)
		sleep *= 2 // экспоненциально увеличиваем паузу
		if e := fn(); e == nil {
			return nil
		} else {
			err = e
		}
	}
	return err
}

// doGet выполняет GET-запрос с retry и обработкой ошибок
func (c *WBClient) doGet(ctx context.Context, url string) ([]byte, error) {
	var respBytes []byte
	err := retry(3, 2*time.Second, func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", c.token)
		req.Header.Set("Accept", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			return fmt.Errorf("http request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf("rate limit 429 from WB API")
		}
		if resp.StatusCode >= 500 {
			return fmt.Errorf("WB API server error: %d", resp.StatusCode)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}
		return nil
	})

	return respBytes, err
}

// parseJSON универсальный парсер []byte → []map[string]interface{}
func parseJSON(data []byte) ([]map[string]interface{}, error) {
	var arr []map[string]interface{}
	var obj map[string]interface{}

	// сначала пробуем как массив
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr, nil
	}

	// если нет — пробуем как объект
	if err := json.Unmarshal(data, &obj); err == nil {
		return []map[string]interface{}{obj}, nil
	}

	return nil, fmt.Errorf("unexpected JSON structure: %s", string(data))
}

//
// ------------------- API Методы -------------------
//

// GetOrders — заказы (исторические)
func (c *WBClient) GetOrders(ctx context.Context, dateFrom, dateTo time.Time) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/orders?dateFrom=%s&flag=0", c.base, dateFrom.Format(time.RFC3339))
	if !dateTo.IsZero() {
		url += "&dateTo=" + dateTo.Format(time.RFC3339)
	}
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseJSON(data)
}

// GetSales — продажи (реализация)
func (c *WBClient) GetSales(ctx context.Context, dateFrom, dateTo time.Time) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/sales?dateFrom=%s", c.base, dateFrom.Format(time.RFC3339))
	if !dateTo.IsZero() {
		url += "&dateTo=" + dateTo.Format(time.RFC3339)
	}
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseJSON(data)
}

// GetStocks — остатки товаров
func (c *WBClient) GetStocks(ctx context.Context, dateFrom time.Time) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/stocks?dateFrom=%s", c.base, dateFrom.Format(time.RFC3339))
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseJSON(data)
}

// GetPrices — текущие цены
func (c *WBClient) GetPrices(ctx context.Context) ([]map[string]interface{}, error) {
	url := "https://common-api.wildberries.ru/api/v1/tariffs" // уточним URL позже
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseJSON(data)
}

// GetTariffs — тарифы доставки
func (c *WBClient) GetTariffs(ctx context.Context) ([]map[string]interface{}, error) {
	url := "https://common-api.wildberries.ru/api/v1/tariffs" // пример, нужно проверить по API
	data, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseJSON(data)
}
