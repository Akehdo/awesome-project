package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const defaultHTTPAddr = ":8080"

type Config struct {
	PostgresDSN string
	HTTPAddr    string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	postgresDSN := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if postgresDSN == "" {
		return nil, errors.New("POSTGRES_DSN is required")
	}

	httpAddr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if httpAddr == "" {
		httpAddr = defaultHTTPAddr
	}

	return &Config{
		PostgresDSN: postgresDSN,
		HTTPAddr:    httpAddr,
	}, nil
}
