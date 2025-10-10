package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// GetOrders godoc
// @Summary Получить заказы из WB API
// @Description Возвращает список заказов за указанный период
// @Tags Orders
// @Param dateFrom query string true "Дата начала (YYYY-MM-DD)"
// @Param dateTo query string false "Дата окончания (YYYY-MM-DD)"
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/orders [get]
func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	dateFrom := r.URL.Query().Get("dateFrom")
	dateTo := r.URL.Query().Get("dateTo")
	if dateFrom == "" {
		http.Error(w, "missing required param: dateFrom", http.StatusBadRequest)
		return
	}

	data, err := h.api.GetOrders(ctx, dateFrom, dateTo)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetOrders failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
