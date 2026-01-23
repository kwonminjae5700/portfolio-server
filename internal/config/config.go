package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port string
	ENV  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type SMTPConfig struct {
	Host     string
	Port     string
	From     string
	Password string
}

type RedisConfig struct {
	Host string
	Port string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "portfolio_db"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			ENV:  getEnv("ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key-change-this"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			From:     getEnv("SMTP_FROM", "kwon777071@gmail.com"),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "minio-api.kwon5700.kr"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "F122SHJRLSTH95AVUZYF"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "rMo8Q+ktptivV79PeWzDdqj10KJOQ1ZGNzVrsxZ8"),
			Bucket:    getEnv("MINIO_BUCKET", "kwon5700-blog"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL", true),
		},
	}
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
