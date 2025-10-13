package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// GetTariffs godoc
// @Summary Получить Комиссия по категориям товаров из WB API
// @Description Метод возвращает данные о комиссии WB по родительским категориям товаров согласно модели продаж.
// @Tags Tariffs
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tariffs [get]
func (h *Handler) GetTariffs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	data, err := h.api.GetTariffs(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetTariffs failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetTariffsBox godoc
// @Summary Получить Комиссия по категориям товаров из WB API
// @Description Метод возвращает данные о комиссии WB по родительским категориям товаров согласно модели продаж.
// @Tags Tariffs
// @Param date query string true "Дата (YYYY-MM-DD)"
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tariffs/box [get]
func (h *Handler) GetTariffsBox(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "missing required param: date", http.StatusBadRequest)
		return
	}

	data, err := h.api.GetTariffsBox(ctx, date)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetTariffs failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetTariffsPallet godoc
// @Summary Получить Комиссия по категориям товаров из WB API
// @Description Метод возвращает данные о комиссии WB по родительским категориям товаров согласно модели продаж.
// @Tags Tariffs
// @Param date query string true "Дата (YYYY-MM-DD)"
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tariffs/pallet [get]
func (h *Handler) GetTariffsPallet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "missing required param: date", http.StatusBadRequest)
		return
	}

	data, err := h.api.GetTariffsPallet(ctx, date)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetTariffs failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
