package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"wildberriesapi/internal/config"
)

type Publisher interface {
	Publish(ctx context.Context, key []byte, v any) error
	Close() error
}

type KafkaPublisher struct {
	writer *kafka.Writer
	logger zerolog.Logger
	topic  string
}

// Создаём Kafka Publisher с проверкой подключения и логированием
func NewKafkaPublisher(cfg config.Config) (Publisher, error) {
	if len(cfg.Kafka.Brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers provided in config")
	}

	// Логгер в stdout
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	brokerAddr := cfg.Kafka.Brokers[0]
	topic := cfg.Kafka.Topic

	logger.Info().Msgf("Connecting to Kafka broker: %s", brokerAddr)

	// Проверим соединение перед созданием writer
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka broker %s: %w", brokerAddr, err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kafka controller: %w", err)
	}
	logger.Info().Msgf("Connected to Kafka controller at %s:%d", controller.Host, controller.Port)

	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		WriteTimeout: 10 * time.Second,
	}

	// Отправим тестовое сообщение
	testMsg := kafka.Message{Value: []byte(`{"status":"ok","message":"Kafka test message"}`)}
	err = w.WriteMessages(context.Background(), testMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to send test message to Kafka: %w", err)
	}

	logger.Info().Msg("✅ Kafka test message sent successfully")

	return &KafkaPublisher{
		writer: w,
		logger: logger,
		topic:  topic,
	}, nil
}

func (p *KafkaPublisher) Publish(ctx context.Context, key []byte, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := kafka.Message{
		Key:   key,
		Value: b,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error().Err(err).Msg("❌ failed to publish message")
		return err
	}

	p.logger.Info().Msgf("✅ message published to topic '%s'", p.topic)
	return nil
}

func (p *KafkaPublisher) Close() error {
	p.logger.Info().Msg("Closing Kafka writer...")
	return p.writer.Close()
}
