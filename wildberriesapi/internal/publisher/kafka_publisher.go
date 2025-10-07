package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"wildberriesapi/internal/config"
)

// Publisher — интерфейс для тестируемости и гибкости
type Publisher interface {
	Publish(ctx context.Context, topic string, key []byte, v any) error
	Close() error
}

// KafkaPublisher — реализация Publisher для Kafka
type KafkaPublisher struct {
	writers map[string]*kafka.Writer
	logger  zerolog.Logger
	brokers []string
}

// NewKafkaPublisher — инициализация Kafka Publisher
func NewKafkaPublisher(cfg config.Config) (Publisher, error) {
	if len(cfg.Kafka.Brokers) == 0 {
		return nil, fmt.Errorf("❌ no Kafka brokers provided in config")
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msgf("🔌 Connecting to Kafka brokers: %v", cfg.Kafka.Brokers)

	// Проверяем соединение с первым брокером
	conn, err := kafka.Dial("tcp", cfg.Kafka.Brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka broker %s: %w", cfg.Kafka.Brokers[0], err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kafka controller: %w", err)
	}
	logger.Info().Msgf("✅ Connected to Kafka controller at %s:%d", controller.Host, controller.Port)

	p := &KafkaPublisher{
		writers: make(map[string]*kafka.Writer),
		logger:  logger,
		brokers: cfg.Kafka.Brokers,
	}

	// Graceful shutdown при SIGTERM / SIGINT
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh
		logger.Info().Msg("🧹 Received shutdown signal, closing Kafka writers...")
		p.Close()
		os.Exit(0)
	}()

	// Проверочный тест
	//testWriter := &kafka.Writer{
	//	Addr:     kafka.TCP(cfg.Kafka.Brokers...),
	//	Topic:    "wb.raw",
	//	Balancer: &kafka.LeastBytes{},
	//}
	//err = testWriter.WriteMessages(context.Background(), kafka.Message{
	//	Value: []byte(`{"status":"ok","message":"Kafka test message"}`),
	//})
	//if err != nil {
	//	return nil, fmt.Errorf("failed to send test message: %w", err)
	//}
	//logger.Info().Msg("✅ Kafka test message sent successfully to topic 'wb.raw.test'")
	//_ = testWriter.Close()

	return p, nil
}

// getWriter — возвращает writer для указанного топика (кеширует)
func (p *KafkaPublisher) getWriter(topic string) *kafka.Writer {
	if w, ok := p.writers[topic]; ok {
		return w
	}
	w := &kafka.Writer{
		Addr:         kafka.TCP(p.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		BatchSize:    10,                     // до 10 сообщений в пачке
		BatchTimeout: 500 * time.Millisecond, // минимальная задержка
	}
	p.writers[topic] = w
	return w
}

// Publish — универсальная публикация сообщения
func (p *KafkaPublisher) Publish(ctx context.Context, topic string, key []byte, v any) error {
	writer := p.getWriter(topic)

	b, err := json.Marshal(v)
	if err != nil {
		p.logger.Error().Err(err).Msg("❌ failed to marshal message")
		return err
	}

	msg := kafka.Message{
		Key:   key,
		Value: b,
		Time:  time.Now(),
	}

	if err := writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error().Err(err).Msgf("❌ failed to publish to topic '%s'", topic)
		return err
	}

	p.logger.Info().Msgf("✅ message published to topic '%s'", topic)
	return nil
}

// Close — закрытие всех Kafka writer’ов
func (p *KafkaPublisher) Close() error {
	for topic, w := range p.writers {
		p.logger.Info().Msgf("🛑 Closing Kafka writer for topic '%s'...", topic)
		if err := w.Close(); err != nil {
			p.logger.Error().Err(err).Msgf("failed to close writer for topic %s", topic)
		}
	}
	return nil
}
