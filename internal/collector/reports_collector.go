package collector

import (
	"context"
	"encoding/json"
	"time"
)

func (r *Collector) CollectDailyReports(ctx context.Context) {
	now := time.Now().AddDate(0, 0, -1) // вчера
	begin := now.Format("2006-01-02") + " 00:00:00"
	end := now.Format("2006-01-02") + " 23:59:59"

	select {
	case <-ctx.Done():
		r.Logger.Warn().Msg("context canceled, stop reports collector")
		return
	default:
	}

	r.Logger.Info().Msgf("📊 Collecting NM Reports for supplier=%d")

	detail, err := r.API.GetNMReportDetailYesterday(ctx, begin, end)
	if err != nil {
		r.Logger.Error().Err(err).Msg("failed to collect nm report detail")
		return
	}

	history, err := r.API.GetNMReportHistoryBatched(ctx, []int{}, begin[:10], end[:10])
	if err != nil {
		r.Logger.Error().Err(err).Msg("failed to collect nm report history")
		return
	}

	reportPayload := map[string]any{
		"date":      begin[:10],
		"detail":    detail,
		"history":   history,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	b, _ := json.Marshal(reportPayload)
	key := []byte("nm_reports_" + begin[:10])
	if err := r.Publisher.Publish(ctx, "wb.raw.reports", key, b); err != nil {
		r.Logger.Error().Err(err).Msg("❌ failed to publish nm report to Kafka")
	} else {
		r.Logger.Info().Msg("✅ NM report published to Kafka topic")
	}
}
