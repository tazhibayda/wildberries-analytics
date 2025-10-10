package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// NMReportItem — структура одного элемента отчёта
type NMReportItem map[string]any

// GetNMReportHistoryBatched — аналог get_nm_report_history_batched()
func (c *WBClient) GetNMReportHistoryBatched(ctx context.Context, nmIDs []int, dateFrom, dateTo string) ([]NMReportItem, error) {
	c.Logger.Info().Msgf("📊 Fetching NM Report history for, period %s..%s", dateFrom, dateTo)

	baseURL := WBBaseURLs["analytics"] + "/nm-report/history"
	out := []NMReportItem{}

	chunks := chunkIntSlice(nmIDs, 20)
	for _, batch := range chunks {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
		}

		payload := map[string]any{
			"period": map[string]string{
				"begin": dateFrom,
				"end":   dateTo,
			},
			"timezone":         "Europe/Moscow",
			"aggregationLevel": "day",
			"nmIDs":            batch,
		}

		body, err := c.doRequest(ctx, http.MethodPost, baseURL, c.Tokens[0], payload)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("nm-report/history error for ")
			continue
		}

		var resp struct {
			Data []NMReportItem `json:"data"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal nm history error")
			continue
		}

		out = append(out, resp.Data...)
		time.Sleep(20 * time.Second)
	}

	return out, nil
}

// GetNMReportDetailYesterday — аналог get_nm_report_detail_yesterday()
func (c *WBClient) GetNMReportDetailYesterday(ctx context.Context, begin, end string) ([]NMReportItem, error) {
	c.Logger.Info().Msgf("📄 Fetching NM Report detail for %s..%s", begin, end)

	baseURL := WBBaseURLs["analytics"] + "/nm-report/detail"
	cards := []NMReportItem{}

	page := 1
	for {
		select {
		case <-ctx.Done():
			return cards, ctx.Err()
		default:
		}

		payload := map[string]any{
			"timezone": "Europe/Moscow",
			"period": map[string]string{
				"begin": begin,
				"end":   end,
			},
			"orderBy": map[string]string{
				"field": "ordersSumRub",
				"mode":  "asc",
			},
			"page": page,
		}

		body, err := c.doRequest(ctx, http.MethodPost, baseURL, c.Tokens[0], payload)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("nm-report/detail error page=%d", page)
			time.Sleep(20 * time.Second)
			continue
		}

		var resp struct {
			Data struct {
				Cards      []NMReportItem `json:"cards"`
				IsNextPage bool           `json:"isNextPage"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			c.Logger.Error().Err(err).Msg("unmarshal nm detail error")
			break
		}

		if len(resp.Data.Cards) == 0 {
			break
		}

		cards = append(cards, resp.Data.Cards...)
		page++
		time.Sleep(20 * time.Second)
		if !resp.Data.IsNextPage {
			break
		}
	}

	return cards, nil
}

// Helpers
func chunkIntSlice(s []int, n int) [][]int {
	if len(s) == 0 {
		return nil
	}
	var chunks [][]int
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
