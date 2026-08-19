package config

import (
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config interface {
	GetServerAddress() string
	GetServerWriteTimeout() time.Duration
	GetServerIdleTimeout() time.Duration
	GetServerReadTimeout() time.Duration
	GetServerReadHeaderTimeout() time.Duration
	GetTemplateDir() string
	GetDomain() string
	GetDatabase() string
	GetLogLevel() string
	GetLogType() string
}

type EnvConfig struct{}

var _ Config = (*EnvConfig)(nil)

func New() (Config, error) {
	_ = godotenv.Load("configs/.env")
	return &EnvConfig{}, nil
}

func (c *EnvConfig) GetServerAddress() string {
	return net.JoinHostPort(
		getEnv("SERVER_LISTEN", "127.0.0.1"),
		getEnv("SERVER_PORT", "1000"),
	)
}

func (c *EnvConfig) GetServerWriteTimeout() time.Duration {
	return getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second)
}

func (c *EnvConfig) GetServerIdleTimeout() time.Duration {
	return getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second)
}

func (c *EnvConfig) GetServerReadTimeout() time.Duration {
	return getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second)
}

func (c *EnvConfig) GetServerReadHeaderTimeout() time.Duration {
	return getEnvDuration("SERVER_READ_HEADER_TIMEOUT", 5*time.Second)
}

func (c *EnvConfig) GetTemplateDir() string {
	return getEnv("TEMPLATE_DIR", "web")
}

func (c *EnvConfig) GetDomain() string {
	return getEnv("DOMAIN", "exemple.com")
}

func (c *EnvConfig) GetDatabase() string {
	return getEnv("DATABASE", "tick.db")
}

// Логи
func (c *EnvConfig) GetLogLevel() string {
	return getEnv("LOG_LEVEL", "info")
}

func (c *EnvConfig) GetLogType() string {
	return getEnv("LOG_TYPE", "tick.db")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		parsedTime, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue
		}
		return parsedTime
	}
	return defaultValue
}
