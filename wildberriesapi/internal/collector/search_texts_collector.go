package collector

import (
	"context"
	"encoding/json"
	_ "fmt"
	"strconv"
	"time"
)

// CollectAndPublishSearchText — вызывает PostSearchTexts для каждого токена и публикует результат в Kafka.
// payload — сформированный body для POST (например, dateFrom/dateTo, filters и т.д.).
func (sc *Collector) CollectAndPublishSearchText(ctx context.Context, payload map[string]interface{}) {
	if ctx.Err() != nil {
		sc.Logger.Warn().Msg("context cancelled before start")
		return
	}

	tokens := sc.API.Tokens

	for i, token := range tokens {
		// graceful stop if cancelled
		if ctx.Err() != nil {
			sc.Logger.Warn().Msg("context cancelled, stopping search-texts collector")
			return
		}

		var supplierID int

		sc.Logger.Info().Msgf("🔎 Calling search-texts for supplier=%d (token_index=%d)", supplierID, i+1)

		respMap, err := sc.API.PostSearchTexts(ctx, payload, token)
		if err != nil {
			sc.Logger.Error().Err(err).Msgf("❌ PostSearchTexts failed for supplier=%d", supplierID)
			// если это таймаут/сетевая ошибка — предлагаем retry или просто продолжим дальше
			continue
		}

		// Если backend вернул объект с "error" — залогируем и опубликуем как ошибочный результат (по необходимости)
		if _, hasErr := respMap["error"]; hasErr {
			sc.Logger.Error().Msgf("search-texts returned error for supplier=%d: %+v", supplierID, respMap["error"])
			// можно публиковать ошибки в отдельный топик; здесь — публикуем тоже (опционально)
		}

		// Добавляем метаданные supplier_id
		respMap["__supplier_id"] = supplierID
		respMap["__fetched_at"] = time.Now().Format(time.RFC3339)

		// Публикация — используем supplierID как key (строка)
		key := []byte(strconv.Itoa(supplierID))
		// Publisher signature: Publish(ctx, topic, key, v)
		if err := sc.Publisher.Publish(ctx, "wb.raw.searchtexts", key, respMap); err != nil {
			sc.Logger.Error().Err(err).Msgf("failed to publish search-texts for supplier=%d", supplierID)
		} else {
			// логируем размер данных для мониторинга (примерно)
			b, _ := json.Marshal(respMap)
			sc.Logger.Info().Msgf("📤 Published search-texts for supplier=%d (bytes=%d)", supplierID, len(b))
		}

		// Небольшая пауза между токенами, чтобы снизить риск 429
		select {
		case <-ctx.Done():
			return
		case <-time.After(sc.API.RetryDelay):
		}
	}
}
