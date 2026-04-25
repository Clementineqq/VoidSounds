package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"voidsounds/internal/components"
	"voidsounds/internal/repository"
	"voidsounds/internal/service"
)

func main() {
	// Инициализация слоёв
	eventRepo := repository.NewEventRepository()
	eventService := service.NewEventService(eventRepo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Роуты
	r.Get("/", homeHandler)
	r.Get("/events", func(w http.ResponseWriter, r *http.Request) {
		eventsHandler(w, r, eventService)
	})

	log.Println("VoidSounds запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	component := components.Home()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request, svc *service.EventService) {
	events, err := svc.GetAllEvents()
	if err != nil {
		http.Error(w, "Ошибка получения мероприятий", http.StatusInternalServerError)
		return
	}

	component := components.Events(events)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
