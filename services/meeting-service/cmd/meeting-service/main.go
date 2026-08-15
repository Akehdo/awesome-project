package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"meeting-service/internal/config"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/service"
	httptransport "meeting-service/internal/transport/http"
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

	meetingRepository := postgres.NewMeetingRepository(db)
	meetingService := service.NewMeetingService(meetingRepository)
	meetingHandler := httptransport.NewMeetingHandler(meetingService)
	router := httptransport.NewRouter(meetingHandler)

	log.Printf("meeting service is listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatal(err)
	}
}
