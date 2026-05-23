package handler

import (
	"net/http"
	"strconv"
	"time"

	"voidsounds/internal/components"
	"voidsounds/internal/domain"
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

// GET /organizer/events/create - форма создания
func (h *EventHandler) ShowCreateForm(w http.ResponseWriter, r *http.Request) {
	components.OrganizerForm(nil, "Создание мероприятия", "/organizer/events", "POST").Render(r.Context(), w)
}

// GET /organizer/events - список своих мероприятий
func (h *EventHandler) GetOrganizerEvents(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetUserID(r)
	events, err := h.service.GetOrganizerEvents(orgID)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		components.OrganizerList(events).Render(r.Context(), w)
	} else {
		// Для прямого входа оборачиваем в Layout
		components.OrganizerPage(events).Render(r.Context(), w)
	}
}

// POST /organizer/events - создание мероприятия
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02T15:04", dateStr)
	if err != nil {
		components.ErrorMessage("Неверный формат даты (ГГГГ-ММ-ДД ЧЧ:ММ)").Render(r.Context(), w)
		return
	}

	price, _ := strconv.Atoi(r.FormValue("price"))
	available, _ := strconv.Atoi(r.FormValue("available"))

	var cityID *int
	if cid := r.FormValue("city_id"); cid != "" {
		val, _ := strconv.Atoi(cid)
		cityID = &val
	}

	var posterURL *string
	if url := r.FormValue("poster_url"); url != "" {
		posterURL = &url
	}

	event := &domain.Event{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Date:        date,
		Address:     r.FormValue("address"),
		Price:       price,
		Available:   available,
		CityID:      cityID,
		PosterURL:   posterURL,
		OrganizerID: middleware.GetUserID(r),
		Status:      "published",
	}

	if err := h.service.CreateEvent(event); err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	// Успех: редирект на список
	w.Header().Set("HX-Redirect", "/organizer/events")
	w.WriteHeader(http.StatusCreated)
}

// DELETE /organizer/events/{id}
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteEvent(id, middleware.GetUserID(r)); err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	// Возвращаем обновлённый список
	h.GetOrganizerEvents(w, r)
}

// PUT /organizer/events/{id} (можно сделать позже, пока заглушка)
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Редактирование в разработке (добавим по запросу)"))
}

// func (h *EventHandler) ShowCreateForm(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) GetOrganizerEvents(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) { ... }
