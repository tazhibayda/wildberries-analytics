package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "wildberriesapi/docs"
	"wildberriesapi/internal/api"
)

// NewRouter создает HTTP маршруты
func NewRouter(api *api.WBClient, log zerolog.Logger) http.Handler {
	r := chi.NewRouter()

	handler := NewHandler(api, log)
	//orders := NewOrdersHandler(api, log)
	//sales := NewSalesHandler(api, log)
	//stocks := NewStocksHandler(api, log)

	r.Get("/api/orders", handler.GetOrders)
	r.Get("/api/sales", handler.GetSales)
	r.Get("/api/stocks", handler.GetStocks)
	r.Get("/api/incomes", handler.GetIncomes)
	r.Get("/api/tariffs", handler.GetTariffs)
	r.Get("/api/tariffs/box", handler.GetTariffsBox)
	r.Get("/api/tariffs/pallet", handler.GetTariffsPallet)
	r.Get("/api/paid_storage/start", handler.StartPaidStorage)
	r.Get("/api/paid_storage/status", handler.GetPaidStorageStatus)
	r.Get("/api/paid_storage/download", handler.GetPaidStorageDownload)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
