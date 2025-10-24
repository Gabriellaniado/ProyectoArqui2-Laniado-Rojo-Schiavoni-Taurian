package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	RabbitMQ RabbitMQConfig
	Solr     SolrConfig
}

type SolrConfig struct {
	Host string
	Port string
	Core string
}

type RabbitMQConfig struct {
	Username  string
	Password  string
	QueueName string
	Host      string
	Port      string
}

func Load() Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}
	return Config{
		Port: getEnv("PORT", "8081"),
		RabbitMQ: RabbitMQConfig{
			Username:  getEnv("RABBITMQ_USER", "admin"),
			Password:  getEnv("RABBITMQ_PASS", "admin"),
			QueueName: getEnv("RABBITMQ_QUEUE_NAME", "items-news"),
			Host:      getEnv("RABBITMQ_HOST", "localhost"),
			Port:      getEnv("RABBITMQ_PORT", "5672"),
		},
		Solr: SolrConfig{
			Host: getEnv("SOLR_HOST", "localhost"),
			Port: getEnv("SOLR_PORT", "8983"),
			Core: getEnv("SOLR_CORE", "demo"),
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
