package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// PaidStorageTask — результат запуска задачи WB "платное хранение"
type PaidStorageTask struct {
	TokenIdx   int    `json:"token_idx"`
	SupplierID int    `json:"supplier_id"`
	TaskID     string `json:"task_id"`
}

// PaidStorageStatus — структура ответа при проверке статуса
type PaidStorageStatus struct {
	Data struct {
		State   string `json:"state"`
		Percent int    `json:"percent"`
	} `json:"data"`
}

// StartPaidStorage запускает сбор данных о платном хранении
func (c *WBClient) StartPaidStorage(ctx context.Context, dateFrom, dateTo string) ([]PaidStorageTask, error) {
	results := make([]PaidStorageTask, 0)
	params := map[string]string{"dateFrom": dateFrom, "dateTo": dateTo}

	for idx, token := range c.Tokens {
		if token == "" {
			continue
		}

		url := fmt.Sprintf("https://seller-analytics-api.wildberries.ru/api/v1/paid_storage?dateFrom=%s&dateTo=%s", dateFrom, dateTo)
		body, err := c.doRequest(ctx, "GET", url, token, params)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("❌ Failed to start paid_storage ")
			continue
		}

		var resp struct {
			Data struct {
				TaskID string `json:"taskId"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal error in start_paid_storage response")
			continue
		}

		if resp.Data.TaskID != "" {
			results = append(results, PaidStorageTask{
				TokenIdx: idx + 1,
				TaskID:   resp.Data.TaskID,
			})
			c.Logger.Info().Msgf("✅ paid_storage started: token_%d, task_id=%s",
				idx+1, resp.Data.TaskID)
		} else {
			c.Logger.Warn().Msgf("⚠️ Unexpected paid_storage start response for")
		}
	}

	return results, nil
}

// GetPaidStorageStatus проверяет статус задачи по token_idx и task_id
func (c *WBClient) GetPaidStorageStatus(ctx context.Context, taskID string) (*PaidStorageStatus, error) {

	url := fmt.Sprintf("https://seller-analytics-api.wildberries.ru/api/v1/paid_storage/tasks/%s/status", taskID)

	body, err := c.doRequest(ctx, "GET", url, c.Tokens[0], nil)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("❌ Failed to get paid_storage status ( task_id=%s)", taskID)
		return nil, err
	}

	var status PaidStorageStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &status, nil
}

// GetPaidStorageDownload скачивает результат задачи по token_idx и task_id
func (c *WBClient) GetPaidStorageDownload(ctx context.Context, taskID string) ([]map[string]any, error) {

	url := fmt.Sprintf("https://seller-analytics-api.wildberries.ru/api/v1/paid_storage/tasks/%s/download", taskID)

	c.Logger.Info().Msgf("⬇️ Downloading paid_storage report ( task_id=%s)", taskID)

	// увеличенный таймаут, потому что ответ может быть большим
	localCtx, cancel := context.WithTimeout(ctx, 4*time.Minute)
	defer cancel()

	body, err := c.doRequest(localCtx, "GET", url, c.Tokens[0], nil)
	if err != nil {
		c.Logger.Error().Err(err).Msgf("❌ Download error (task_id=%s)", taskID)
		return nil, err
	}

	var data []map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON download: %w", err)
	}

	c.Logger.Info().Msgf("✅ Paid storage report downloaded successfully (records=%d)", len(data))
	return data, nil
}
