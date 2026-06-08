package main

import (
	"CRUD/config"
	"CRUD/internal/handler"
	"CRUD/internal/repository"
	"CRUD/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("/Users/mihailignatev/Desktop/WB I2/l2.18/config/config.yaml")
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

	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal(err)
	}
}
