package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PostSearchTexts выполняет POST-запрос к endpoint'у search-texts.
// Возвращает JSON как map[string]interface{}.
// Если сервер вернул не-200 — возвращаем map{"error": {"status": code, "message": body}} (без error),
// чтобы поведение было похоже на Python-реализацию.
// Если токен == "" — используется первый токен в c.Tokens.
// NOTE: Проверь правильный path для своего окружения; здесь примерный путь.
func (c *WBClient) PostSearchTexts(ctx context.Context, payload map[string]interface{}, token string) (map[string]interface{}, error) {
	// выбор токена
	useToken := token
	if useToken == "" {
		if len(c.Tokens) == 0 {
			return nil, fmt.Errorf("no wb tokens configured")
		}
		useToken = c.Tokens[0]
	}

	// endpoint — проверьте ваш фактический путь: тут примерный
	// (в Python у тебя ENDPOINTS['search_texts'])
	endpoint := WBBaseURLs["analytics"] + "/nm-report/search-texts"

	// marshal payload
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	// настроим client с увеличенным таймаутом (по аналогии с sock_read=240)
	reqTimeout := c.ReadTimeout
	if reqTimeout == 0 {
		reqTimeout = 240 * time.Second
	}

	httpClient := &http.Client{
		Timeout: reqTimeout,
	}

	// создаём контекст с таймаутом на случай, если caller дал общий ctx без deadline
	ctxReq, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxReq, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", useToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		// таймаут/сетевая ошибка — возвращаем её как ошибку (как в Python: Timeout -> exception)
		return nil, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
	}

	// Если статус != 200 — возвращаем map с ключом error (поведение Python)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := string(respBytes)
		out := map[string]interface{}{
			"error": map[string]interface{}{
				"status":  resp.StatusCode,
				"message": msg,
			},
		}
		return out, nil
	}

	// Парсим JSON минимально: если это object -> возвращаем его как map[string]interface{}
	// если это array -> упакуем в {"data": [...]}
	var parsed any
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		// Попытка fallback: вернуть текст ошибки в структуре
		out := map[string]interface{}{
			"error": map[string]interface{}{
				"status":  "parse_error",
				"message": fmt.Sprintf("failed to parse json response: %v; raw: %s", err, string(respBytes)),
			},
		}
		return out, nil
	}

	switch v := parsed.(type) {
	case map[string]interface{}:
		return v, nil
	default:
		// например []any -> упакуем в поле data, чтобы результат всегда map[string]interface{}
		out := map[string]interface{}{
			"data": v,
		}
		return out, nil
	}
}
