package handlers

import (
	"github.com/rs/zerolog"
	"wildberriesapi/internal/api"
)

type Handler struct {
	api    *api.WBClient
	logger zerolog.Logger
}

func NewHandler(api *api.WBClient, logger zerolog.Logger) *Handler {
	return &Handler{
		api:    api,
		logger: logger,
	}
}
