package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const defaultHTTPAddr = ":8080"

type Config struct {
	PostgresDSN       string
	HTTPAddr          string
	MinioEndpoint     string // Адрес конечной точки Minio
	MinioBucketName   string // Название конкретного бакета в Minio
	MinioRootUser     string // Имя пользователя для доступа к Minio
	MinioRootPassword string // Пароль для доступа к Minio
	MinioUseSSL       bool   // Переменная, отвечающая за
}

func NewConfig() (*Config, error) {
	if err := loadDotEnv(); err != nil {
		return nil, err
	}

	postgresDSN := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if postgresDSN == "" {
		return nil, errors.New("POSTGRES_DSN is required")
	}

	httpAddr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if httpAddr == "" {
		httpAddr = defaultHTTPAddr
	}

	minioEndpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	if minioEndpoint == "" {
		return nil, errors.New("MINIO_ENDPOINT is required")
	}

	minioBucketName := strings.TrimSpace(os.Getenv("MINIO_BUCKET_NAME"))
	if minioBucketName == "" {
		return nil, errors.New("MINIO_BUCKET_NAME is required")
	}

	minioRootUser := strings.TrimSpace(os.Getenv("MINIO_ROOT_USER"))
	if minioRootUser == "" {
		return nil, errors.New("MINIO_ROOT_USER is required")
	}

	minioRootPassword := strings.TrimSpace(os.Getenv("MINIO_ROOT_PASSWORD"))
	if minioRootPassword == "" {
		return nil, errors.New("MINIO_ROOT_PASSWORD is required")
	}

	minioUseSSLRaw := strings.TrimSpace(os.Getenv("MINIO_USE_SSL"))
	if minioUseSSLRaw == "" {
		return nil, errors.New("MINIO_USE_SSL is required")
	}

	minioUseSSL, err := strconv.ParseBool(minioUseSSLRaw)
	if err != nil {
		return nil, fmt.Errorf(
			"MINIO_USE_SSL must be true or false: %w",
			err,
		)
	}

	return &Config{
		PostgresDSN:       postgresDSN,
		HTTPAddr:          httpAddr,
		MinioEndpoint:     minioEndpoint,
		MinioBucketName:   minioBucketName,
		MinioRootUser:     minioRootUser,
		MinioRootPassword: minioRootPassword,
		MinioUseSSL:       minioUseSSL,
	}, nil
}

func loadDotEnv() error {
	paths := []string{
		".env",
		filepath.Join("..", "..", ".env"),
	}

	for _, path := range paths {
		err := godotenv.Load(path)
		if err == nil {
			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("load %s: %w", path, err)
		}
	}

	return nil
}
