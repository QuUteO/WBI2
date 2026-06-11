package main

import (
	"CRUD/config"
	"CRUD/internal/handler"
	"CRUD/internal/middleware"
	"CRUD/internal/repository"
	"CRUD/internal/service"
	"flag"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository()
	srv := service.NewService(repo)
	calendarHandler := handler.NewHandler(srv)

	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", calendarHandler.CreateEvent)
	mux.HandleFunc("/update_event", calendarHandler.UpdateEvent)
	mux.HandleFunc("/delete_event", calendarHandler.DeleteEvent)
	mux.HandleFunc("/events_for_day", calendarHandler.GetEventsForDay)
	mux.HandleFunc("/events_for_week", calendarHandler.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", calendarHandler.GetEventsForMonth)

	loggedMux := middleware.LoggingMiddleware(mux)

	if err := http.ListenAndServe(cfg.Addr, loggedMux); err != nil {
		log.Fatal(err)
	}
}
