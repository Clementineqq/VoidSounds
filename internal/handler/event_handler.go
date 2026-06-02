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
	"github.com/skip2/go-qrcode"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) Home(w http.ResponseWriter, r *http.Request) {
	components.Home().Render(r.Context(), w)
}

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
		log.Printf("❌ Ошибка фильтрации: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cities, _ := h.service.GetAllCities()
	genres, _ := h.service.GetAllGenres()

	if r.Header.Get("HX-Request") == "true" {
		components.EventsContent(events, cities, genres, city, genre, search).Render(r.Context(), w)
	} else {
		components.Events(events, cities, genres, city, genre, search).Render(r.Context(), w)
	}
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	event, err := h.service.GetEventWithGenres(id)
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

func (h *EventHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID мероприятия").Render(r.Context(), w)
		return
	}

	err = h.service.BuyTicket(eventID, userID)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	components.TicketSuccess(eventID).Render(r.Context(), w)
}

func (h *EventHandler) ShowCreateForm(w http.ResponseWriter, r *http.Request) {
	genres, _ := h.service.GetAllGenres()
	components.OrganizerForm(nil, "Создание мероприятия", "/organizer/events", "POST", genres).Render(r.Context(), w)
}

func (h *EventHandler) ShowEditForm(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID").Render(r.Context(), w)
		return
	}

	event, err := h.service.GetEventByIDForEdit(id)
	if err != nil {
		components.ErrorMessage("Мероприятие не найдено").Render(r.Context(), w)
		return
	}

	if event.OrganizerID != middleware.GetUserID(r) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Доступ запрещён"))
		return
	}

	event.Genres, _ = h.service.GetGenresByEventID(id)
	genres, _ := h.service.GetAllGenres()
	components.OrganizerForm(event, "Редактирование мероприятия", "/organizer/events/"+idStr+"/update", "POST", genres).Render(r.Context(), w)
}

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
		components.OrganizerPage(events).Render(r.Context(), w)
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 << 20)

	dateStr := r.FormValue("date")
	date, err := time.ParseInLocation("2006-01-02T15:04", dateStr, time.Local)
	if err != nil {
		components.ErrorMessage("Неверный формат даты").Render(r.Context(), w)
		return
	}

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

	var posterURL *string
	file, header, err := r.FormFile("poster")
	if err == nil && file != nil {
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
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
			url := "/static/uploads/" + strings.ReplaceAll(filename, "\\", "/")
			posterURL = &url
			log.Printf("🖼️ Постер сохранён: %s", url)
		}
	}

	// Обработка жанров
	var genreIDs []int
	for _, idStr := range r.Form["genres"] {
		if id, err := strconv.Atoi(idStr); err == nil {
			genreIDs = append(genreIDs, id)
		}
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

	if err := h.service.CreateEventWithGenres(event, genreIDs); err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/organizer/events")
	w.WriteHeader(http.StatusCreated)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(5 << 20)

	dateStr := r.FormValue("date")
	date, err := time.ParseInLocation("2006-01-02T15:04", dateStr, time.Local)
	if err != nil {
		components.ErrorMessage("Неверный формат даты").Render(r.Context(), w)
		return
	}

	// Обработка типа мероприятия (платное/бесплатное)
	eventType := r.FormValue("event_type")
	price := 0
	if eventType != "free" {
		p, _ := strconv.Atoi(r.FormValue("price"))
		price = p
	}

	available, _ := strconv.Atoi(r.FormValue("available"))

	var posterURL *string
	file, header, err := r.FormFile("poster")
	if err == nil && file != nil {
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
			os.MkdirAll("static/uploads", 0755)
			safeName := strings.ReplaceAll(r.FormValue("title"), " ", "_")
			filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), safeName, ext)
			path := filepath.Join("static/uploads", filename)
			out, err := os.Create(path)
			if err == nil {
				defer out.Close()
				io.Copy(out, file)
				url := "/static/uploads/" + strings.ReplaceAll(filename, "\\", "/")
				posterURL = &url
			}
		}
	}

	// Обработка жанров
	var genreIDs []int
	for _, idStr := range r.Form["genres"] {
		if id, err := strconv.Atoi(idStr); err == nil {
			genreIDs = append(genreIDs, id)
		}
	}

	// Получаем статус из формы, если он есть
	status := r.FormValue("status")
	if status == "" {
		status = "published" // Значение по умолчанию
	}

	event := &domain.Event{
		ID:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Date:        date,
		Address:     r.FormValue("address"),
		Price:       price,
		Available:   available,
		PosterURL:   posterURL,
		Status:      status,
	}

	err = h.service.UpdateEventWithGenres(id, middleware.GetUserID(r), event, genreIDs)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/organizer/events")
	w.WriteHeader(http.StatusOK)
}

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

	h.GetOrganizerEvents(w, r)
}

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

	h.GetOrganizerEvents(w, r)
}

func (h *EventHandler) TicketQR(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	qrData := fmt.Sprintf("https://voidsounds.ru/ticket/verify/%s", ticketID)
	png, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Ошибка генерации QR", 500)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}
