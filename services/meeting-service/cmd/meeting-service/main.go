package main

import (
	"context"
	"database/sql"
	"log"
	"meeting-service/storage"
	"net/http"
	"time"

	"meeting-service/internal/config"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/service"
	httptransport "meeting-service/internal/transport/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("load config: ", err)
	}

	db, err := sql.Open("pgx", cfg.PostgresDSN)
	if err != nil {
		log.Fatal("open postgres connection: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("connect to postgres: ", err)
	}

	minioClient, err := minio.New(
		cfg.MinioEndpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				cfg.MinioRootUser,
				cfg.MinioRootPassword,
				"",
			),
			Secure: cfg.MinioUseSSL,
		},
	)
	if err != nil {
		log.Fatal("create MinIO client: ", err)
	}

	meetingStorage := storage.NewStorage(
		minioClient,
		cfg.MinioBucketName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, cfg.MinioBucketName)
	if err != nil {
		log.Fatal("check MinIO bucket: ", err)
	}
	if !exists {
		log.Fatalf("MinIO bucket %q does not exist", cfg.MinioBucketName)
	}

	meetingRepository := postgres.NewMeetingRepository(db)
	meetingService := service.NewMeetingService(meetingRepository, meetingStorage)
	meetingHandler := httptransport.NewMeetingHandler(meetingService)
	router := httptransport.NewRouter(meetingHandler)

	log.Printf("meeting service is listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}
