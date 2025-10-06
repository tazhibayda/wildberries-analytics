package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type Config struct {
	WBToken      string
	PollInterval time.Duration
	Kafka        KafkaConfig
	LogLevel     string
	HTTPTimeout  time.Duration
}

func Load() Config {
	v := viper.New()
	v.SetEnvPrefix("WB")
	v.AutomaticEnv()

	v.SetDefault("POLL_INTERVAL", "30m")
	v.SetDefault("KAFKA_TOPIC", "wb.raw")
	v.SetDefault("KAFKA_BROKERS", "kafka:9092")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("HTTP_TIMEOUT", "30s")

	poll, _ := time.ParseDuration(v.GetString("POLL_INTERVAL"))
	httpTimeout, _ := time.ParseDuration(v.GetString("HTTP_TIMEOUT"))

	brokers := []string{}
	rawBrokers := v.GetString("KAFKA_BROKERS")
	for _, b := range splitAndTrim(rawBrokers, ",") {
		if b != "" {
			brokers = append(brokers, b)
		}
	}

	return Config{
		WBToken:      v.GetString("TOKEN"),
		PollInterval: poll,
		Kafka: KafkaConfig{
			Brokers: brokers,
			Topic:   v.GetString("KAFKA_TOPIC"),
		},
		LogLevel:    v.GetString("LOG_LEVEL"),
		HTTPTimeout: httpTimeout,
	}
}

func splitAndTrim(s, sep string) []string {
	out := []string{}
	for _, p := range strings.Split(s, sep) {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}
