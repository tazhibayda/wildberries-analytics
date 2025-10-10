package collector

import (
	"context"
	"time"

	"wildberriesapi/internal/api"
	"wildberriesapi/internal/config"
	"wildberriesapi/internal/publisher"

	"github.com/rs/zerolog"
)

type AdvertsCollector struct {
	cfg       config.Config
	api       *api.WBClient
	publisher publisher.Publisher
	logger    zerolog.Logger
}

func NewAdvertsCollector(cfg config.Config, client *api.WBClient, pub publisher.Publisher, log zerolog.Logger) *AdvertsCollector {
	return &AdvertsCollector{
		cfg:       cfg,
		api:       client,
		publisher: pub,
		logger:    log,
	}
}

func (c *AdvertsCollector) Run(ctx context.Context) {
	c.logger.Info().Msg("🚀 Starting AdvertsCollector loop")

	ticker := time.NewTicker(c.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info().Msg("🛑 AdvertsCollector stopped")
			return
		case <-ticker.C:
			c.collectAndPublish(ctx)
		}
	}
}

func (c *AdvertsCollector) collectAndPublish(ctx context.Context) {
	c.logger.Info().Msg("📥 Collecting adverts from WB API")

	adverts, err := c.api.GetAdverts(ctx)
	if err != nil {
		c.logger.Error().Err(err).Msg("❌ failed to fetch adverts")
		return
	}

	count := 0
	for _, adv := range adverts {
		err := c.publisher.Publish(ctx, "wb.raw.adverts", nil, adv)
		if err == nil {
			count++
		}
	}

	c.logger.Info().Msgf("✅ Published %d advert records to Kafka topic 'wb.adverts'", count)
}
