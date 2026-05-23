package handler

import (
	"net/http"
	"strconv"

	"voidsounds/internal/components"
	"voidsounds/internal/middleware"
	"voidsounds/internal/service"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// Главная
func (h *EventHandler) Home(w http.ResponseWriter, r *http.Request) {
	components.Home().Render(r.Context(), w)
}

// Список мероприятий
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		components.EventsContent(events).Render(r.Context(), w)
	} else {
		components.Events(events).Render(r.Context(), w)
	}
}

// Детальная страница события
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	event, err := h.service.GetEventByID(id)
	if err != nil {
		http.Error(w, "Событие не найдено", 404)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		components.EventDetailContent(event).Render(r.Context(), w)
	} else {
		components.EventDetailPage(event).Render(r.Context(), w)
	}
}

// POST /event/{id}/buy - покупка билета (HTMX)
func (h *EventHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	userID := middleware.GetUserID(r)
	if userID == 0 {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Получаем ID мероприятия из URL
	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID мероприятия").Render(r.Context(), w)
		return
	}

	// Пытаемся купить билет
	err = h.service.BuyTicket(eventID, userID)
	if err != nil {
		// Показываем ошибку
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	// Успех! Показываем компонент подтверждения
	// (создадим его ниже)
	components.TicketSuccess(eventID).Render(r.Context(), w)
}

// func (h *EventHandler) ShowCreateForm(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) GetOrganizerEvents(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) { ... }
