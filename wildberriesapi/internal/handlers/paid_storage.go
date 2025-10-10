package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// StartPaidStorage godoc
// @Summary Создать отчёт из WB API
// @Description Метод создаёт задание на генерацию отчёта о платном хранении.
//
//	Можно получить отчёт максимум за 8 дней.
//
// @Tags Paid Storage
// @Param dateFrom query string true "Дата начала (YYYY-MM-DD)"
// @Param dateTo query string false "Дата окончания (YYYY-MM-DD)"
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/paid_storage/start [get]
func (h *Handler) StartPaidStorage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	dateFrom := r.URL.Query().Get("dateFrom")
	dateTo := r.URL.Query().Get("dateTo")
	if dateFrom == "" {
		http.Error(w, "missing required param: dateFrom", http.StatusBadRequest)
		return
	}
	if dateTo == "" {
		http.Error(w, "missing required param: dateFrom", http.StatusBadRequest)
		return
	}

	data, err := h.api.StartPaidStorage(ctx, dateFrom, dateTo)
	if err != nil {
		h.logger.Error().Err(err).Msg("StartPaidStorage failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPaidStorageStatus godoc
// @Summary Проверить статус из WB API
// @Description Возвращает статус задания на генерацию отчёта о платном хранении заказов за указанный период
// @Tags Paid Storage
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/paid_storage/storage [get]
func (h *Handler) GetPaidStorageStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	taskId := r.URL.Query().Get("taskId")
	if taskId == "" {
		h.logger.Error().Msg("missing required param: taskId")
		http.Error(w, "missing required param: taskId", http.StatusBadRequest)
		return
	}

	data, err := h.api.GetPaidStorageStatus(ctx, taskId)
	if err != nil {
		h.logger.Error().Err(err).Msg("StartPaidStorage failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetPaidStorageDownload godoc
// @Summary Получить отчёт из WB API
// @Description Метод возвращает отчёт о платном хранении по ID задания на генерацию.
// @Tags Paid Storage
// @Param dateFrom query string true "Дата начала (YYYY-MM-DD)"
// @Param dateTo query string false "Дата окончания (YYYY-MM-DD)"
// @Success 200 {object} []map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/paid_storage/download [get]
func (h *Handler) GetPaidStorageDownload(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	taskId := r.URL.Query().Get("taskId")
	if taskId == "" {
		h.logger.Error().Msg("missing required param: taskId")
		http.Error(w, "missing required param: taskId", http.StatusBadRequest)
		return
	}

	data, err := h.api.GetPaidStorageDownload(ctx, taskId)
	if err != nil {
		h.logger.Error().Err(err).Msg("StartPaidStorage failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
