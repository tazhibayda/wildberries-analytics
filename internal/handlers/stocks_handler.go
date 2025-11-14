package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// GetStocks godoc
// @Summary Получить остатки из WB API
// @Description Возвращает текущие остатки
// @Tags Stocks
// @Param dateFrom query string true "Дата начала (YYYY-MM-DD)"
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/stocks [get]
func (h *Handler) GetStocks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	dateFrom := r.URL.Query().Get("dateFrom")
	if dateFrom == "" {
		http.Error(w, "missing required param: dateFrom", http.StatusBadRequest)
		return
	}

	data, err := h.api.GetStocks(ctx, dateFrom)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetStocks failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
