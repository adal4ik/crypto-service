package config

import (
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	AppPort  string
	LogLevel string
}

type PostgresConfig struct {
	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
	PostgresDSN string
}

type CollectorConfig struct {
	Interval   time.Duration
	ApiBaseURL string
}

type Config struct {
	App       AppConfig
	Postgres  PostgresConfig
	Collector CollectorConfig
}

func LoadConfig() *Config {
	collectorIntervalSec, err := strconv.Atoi(getEnv("COLLECTOR_INTERVAL_SECONDS", "60"))
	if err != nil {
		// Если в .env указано не число, ставим значение по умолчанию.
		collectorIntervalSec = 5
	}
	cfg := &Config{
		App: AppConfig{
			AppPort:  getEnv("APP_PORT", "8080"),
			LogLevel: getEnv("LOG_LEVEL", "DEBUG"),
		},
		Postgres: PostgresConfig{
			DBHost:      getEnv("DB_HOST", "db"),
			DBPort:      getEnv("DB_PORT", "5432"),
			DBName:      getEnv("DB_NAME", "asd"),
			DBUser:      getEnv("DB_USER", "postgres"),
			DBPassword:  getEnv("DB_PASSWORD", "supersecret"),
			PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:supersecret@db:5432/asd?sslmode=disable"),
		},
		Collector: CollectorConfig{
			Interval:   time.Duration(collectorIntervalSec) * time.Second,
			ApiBaseURL: getEnv("COINGECKO_API_URL", "https://api.coingecko.com/api/v3/simple/price"),
		},
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
