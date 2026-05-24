package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// Список мероприятий (с поддержкой фильтров)
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	genre := r.URL.Query().Get("genre")
	search := r.URL.Query().Get("search")

	var events domain.Events
	var err error

	if city != "" || genre != "" || search != "" {
		events, err = h.service.GetEventsWithFilters(city, genre, search)
	} else {
		events, err = h.service.GetAllEvents()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем города и жанры для фильтров
	cities, _ := h.service.GetAllCities()
	genres, _ := h.service.GetAllGenres()

	if r.Header.Get("HX-Request") == "true" {
		components.EventsContent(events, cities, genres, city, genre, search).Render(r.Context(), w)
	} else {
		components.Events(events, cities, genres, city, genre, search).Render(r.Context(), w)
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
	r.ParseMultipartForm(5 << 20) // лимит 5MB

	dateStr := r.FormValue("date")
	date, err := time.ParseInLocation("2006-01-02T15:04", dateStr, time.Local)
	if err != nil {
		components.ErrorMessage("Неверный формат даты").Render(r.Context(), w)
		return
	}

	// Парсинг цены / свободного входа
	isFree := r.FormValue("is_free") == "on"
	price := 0
	if !isFree {
		p, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			components.ErrorMessage("Неверная цена").Render(r.Context(), w)
			return
		}
		price = p
	}
	available, _ := strconv.Atoi(r.FormValue("available"))

	var cityID *int
	if cid := r.FormValue("city_id"); cid != "" {
		val, _ := strconv.Atoi(cid)
		cityID = &val
	}

	// === ЗАГРУЗКА ПОСТЕРА ===
	var posterURL *string
	file, handler, err := r.FormFile("poster")
	if err == nil && file != nil {
		defer file.Close()

		// Разрешаем только картинки
		ext := strings.ToLower(filepath.Ext(handler.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			components.ErrorMessage("Разрешены только JPG, PNG, WEBP").Render(r.Context(), w)
			return
		}

		os.MkdirAll("static/uploads", 0755)
		safeName := strings.ReplaceAll(r.FormValue("title"), " ", "_")
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), safeName, ext)
		path := filepath.Join("static/uploads", filename)

		out, err := os.Create(path)
		if err == nil {
			defer out.Close()
			io.Copy(out, file)
			// Для Windows заменяем \ на /, чтобы браузер открывал
			url := "/static/uploads/" + strings.ReplaceAll(filename, "\\", "/")
			posterURL = &url
			log.Printf("🖼️ Постер сохранён: %s", url)
		}
	}
	// =======================

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

// GET /profile - личный кабинет (история билетов)
func (h *EventHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tickets, err := h.service.GetUserTickets(userID)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		components.ProfileContent(tickets).Render(r.Context(), w)
	} else {
		components.ProfilePage(tickets).Render(r.Context(), w)
	}
}

// POST /organizer/events/{id}/status - смена статуса
func (h *EventHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newStatus := r.FormValue("status")
	if newStatus == "" {
		components.ErrorMessage("Выберите статус").Render(r.Context(), w)
		return
	}

	// Разрешаем только валидные статусы
	allowed := map[string]bool{"published": true, "draft": true, "cancelled": true}
	if !allowed[newStatus] {
		components.ErrorMessage("Недопустимый статус").Render(r.Context(), w)
		return
	}

	err = h.service.UpdateStatus(id, middleware.GetUserID(r), newStatus)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	// Возвращаем обновлённый список
	h.GetOrganizerEvents(w, r)
}

// func (h *EventHandler) ShowCreateForm(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) GetOrganizerEvents(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) { ... }
// func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) { ... }
