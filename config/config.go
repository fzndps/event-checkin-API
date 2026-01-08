// Package config untuk menyimpan semua konfigurasi aplikasi
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	App      AppConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type AppConfig struct {
	Port string
}

type JWTConfig struct {
	Secret string
	Expiry int
}

type SMTPConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func LoadConfig(filename ...string) (*Config, error) {
	if err := godotenv.Load(filename...); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},

		App: AppConfig{
			Port: os.Getenv("SERVER_PORT"),
		},

		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Expiry: getEnvAsInt("JWT_EXPIRY", 72),
		},

		SMTP: SMTPConfig{
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 578),
			SMTPUsername: os.Getenv("SMTP_USERNAME"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
			SMTPFrom:     os.Getenv("SMTP_FROM"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	// Cek database config
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Cek JWT secret (sangat penting untuk keamanan!)
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	// Cek app port
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	// Cek email config (jika ingin kirim email)
	if c.SMTP.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if c.SMTP.SMTPUsername == "" {
		return fmt.Errorf("SMTP_USERNAME is required")
	}
	if c.SMTP.SMTPPassword == "" {
		return fmt.Errorf("SMTP_PASSWORD is required")
	}

	// Semua validasi passed
	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

func getEnvAsInt(key string, defaultValue int) int {
	valueSTR := os.Getenv(key)

	if valueSTR == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueSTR)
	if err != nil {
		return defaultValue
	}

	return value
}
