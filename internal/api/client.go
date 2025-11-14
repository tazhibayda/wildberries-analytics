package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

// WBClient — клиент для Wildberries Client
type WBClient struct {
	BaseURL map[string]string
	Tokens  []string
	//SupplierIDs []int
	Client      *http.Client
	Logger      zerolog.Logger
	RetryDelay  time.Duration
	MaxRetries  int
	ReadTimeout time.Duration
}

// NewWBClient создаёт новый WB Client клиент
func NewWBClient(tokens []string, log zerolog.Logger) *WBClient {
	return &WBClient{
		BaseURL:     WBBaseURLs,
		Tokens:      tokens,
		Client:      &http.Client{Timeout: 60 * time.Second},
		Logger:      log,
		RetryDelay:  2 * time.Second,
		MaxRetries:  5,
		ReadTimeout: 120 * time.Second,
	}
}

// doRequest выполняет GET-запрос с retry и обработкой ошибок
func (c *WBClient) doRequest(ctx context.Context, method, url, token string, payload any) ([]byte, error) {
	const maxJSONSize = 20 << 20 // 20 MB
	var body io.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	var resp *http.Response
	for attempt := 1; attempt <= c.MaxRetries; attempt++ {
		resp, err = c.Client.Do(req)
		if err != nil {
			if attempt < c.MaxRetries {
				time.Sleep(c.RetryDelay * time.Duration(attempt))
				continue
			}
			return nil, fmt.Errorf("network error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 && attempt < c.MaxRetries {
			c.Logger.Warn().Msgf("WB Client %d — retry %d/%d", resp.StatusCode, attempt, c.MaxRetries)
			time.Sleep(c.RetryDelay * time.Duration(attempt))
			continue
		}

		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("unauthorized (401)")
		}

		if resp.StatusCode == 429 {
			c.Logger.Warn().Msg("Too many requests (429), waiting 60s...")
			time.Sleep(60 * time.Second)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(io.LimitReader(resp.Body, maxJSONSize))
			return nil, fmt.Errorf("unexpected WB Client status %d: %s", resp.StatusCode, string(b))
		}

		b, err := io.ReadAll(io.LimitReader(resp.Body, maxJSONSize))
		if err != nil {
			return nil, fmt.Errorf("read body error: %w", err)
		}

		return b, nil
	}

	return nil, errors.New("max retries reached")
}
