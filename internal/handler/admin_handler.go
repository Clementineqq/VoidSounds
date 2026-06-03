package handler

import (
	"net/http"
	"strconv"
	"time"

	"voidsounds/internal/components"
	"voidsounds/internal/domain"
	"voidsounds/internal/service"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	eventService *service.EventService
	userService  *service.UserService
}

func NewAdminHandler(eventService *service.EventService, userService *service.UserService) *AdminHandler {
	return &AdminHandler{
		eventService: eventService,
		userService:  userService,
	}
}

// GET /admin
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	users, _ := h.userService.GetAllUsers()
	events, _ := h.eventService.GetAllEventsForAdmin()

	stats := components.AdminStats{
		TotalUsers:   len(users),
		TotalEvents:  len(events),
		TotalTickets: h.calculateTotalTickets(events),
	}

	if r.Header.Get("HX-Request") == "true" {
		components.AdminDashboardContent(stats, users, events).Render(r.Context(), w)
	} else {
		components.AdminDashboard(stats, users, events).Render(r.Context(), w)
	}
}

// GET /admin/users - список пользователей
func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		components.ErrorMessage("Ошибка загрузки пользователей").Render(r.Context(), w)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		components.AdminUsersList(users).Render(r.Context(), w)
	} else {
		// Для обычного запроса оборачиваем в Layout
		components.AdminUsersPage(users).Render(r.Context(), w)
	}
}

// GET /admin/users/{id} - информация о пользователе
func (h *AdminHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID").Render(r.Context(), w)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		components.ErrorMessage("Пользователь не найден").Render(r.Context(), w)
		return
	}

	// Получаем статистику пользователя
	userEvents, _ := h.eventService.GetEventsByOrganizer(id)
	userTickets, _ := h.eventService.GetUserTickets(id)

	component := components.AdminUserDetail(user, userEvents, userTickets)
	component.Render(r.Context(), w)
}

// POST /admin/users/{id}/role - смена роли пользователя
func (h *AdminHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID").Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		components.ErrorMessage("Ошибка обработки формы").Render(r.Context(), w)
		return
	}

	newRole := r.FormValue("role")
	allowed := map[string]bool{"viewer": true, "organizer": true, "admin": true}
	if !allowed[newRole] {
		components.ErrorMessage("Недопустимая роль").Render(r.Context(), w)
		return
	}

	err = h.userService.ChangeUserRole(id, newRole)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/admin/users")
	w.WriteHeader(http.StatusOK)
}

// POST /admin/users/{id}/ban - бан/разбан пользователя
func (h *AdminHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID").Render(r.Context(), w)
		return
	}

	err = h.userService.BanUser(id)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/admin/users")
	w.WriteHeader(http.StatusOK)
}

// GET /admin/events - список всех мероприятий
func (h *AdminHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.GetAllEventsForAdmin()
	if err != nil {
		components.ErrorMessage("Ошибка загрузки мероприятий").Render(r.Context(), w)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		components.AdminEventsList(events).Render(r.Context(), w)
	} else {
		components.AdminEventsPage(events).Render(r.Context(), w)
	}
}

// POST /admin/events/{id}/status - смена статуса мероприятия
func (h *AdminHandler) ChangeEventStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID").Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		components.ErrorMessage("Ошибка обработки формы").Render(r.Context(), w)
		return
	}

	newStatus := r.FormValue("status")
	allowed := map[string]bool{"published": true, "draft": true, "cancelled": true}
	if !allowed[newStatus] {
		components.ErrorMessage("Недопустимый статус").Render(r.Context(), w)
		return
	}

	// Получаем ID организатора для проверки
	event, err := h.eventService.GetEventByID(id)
	if err != nil {
		components.ErrorMessage("Мероприятие не найдено").Render(r.Context(), w)
		return
	}

	err = h.eventService.UpdateStatus(id, event.OrganizerID, newStatus)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/admin/events")
	w.WriteHeader(http.StatusOK)
}

// DELETE /admin/events/{id} - удаление мероприятия
func (h *AdminHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.eventService.DeleteEventAdmin(id); err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/admin/events")
	w.WriteHeader(http.StatusOK)
}

// Вспомогательные методы
func (h *AdminHandler) calculateTotalTickets(events domain.Events) int {
	total := 0
	for _, e := range events {
		total += e.Available
	}
	return total
}

// GET /admin/events/{id}/edit - форма редактирования
func (h *AdminHandler) EditEventForm(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID").Render(r.Context(), w)
		return
	}

	event, err := h.eventService.GetEventByIDForEdit(id)
	if err != nil {
		components.ErrorMessage("Мероприятие не найдено").Render(r.Context(), w)
		return
	}

	event.Genres, _ = h.eventService.GetGenresByEventID(id)
	genres, _ := h.eventService.GetAllGenres()

	component := components.AdminEditEventForm(event, genres)
	component.Render(r.Context(), w)
}

// POST /admin/events/{id}/update - обновление мероприятия
func (h *AdminHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		components.ErrorMessage("Неверный ID").Render(r.Context(), w)
		return
	}

	r.ParseMultipartForm(5 << 20)

	dateStr := r.FormValue("date")
	date, err := time.ParseInLocation("2006-01-02T15:04", dateStr, time.Local)
	if err != nil {
		components.ErrorMessage("Неверный формат даты").Render(r.Context(), w)
		return
	}

	price, _ := strconv.Atoi(r.FormValue("price"))
	available, _ := strconv.Atoi(r.FormValue("available"))

	event := &domain.Event{
		ID:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Date:        date,
		Address:     r.FormValue("address"),
		Price:       price,
		Available:   available,
		Status:      r.FormValue("status"),
	}

	err = h.eventService.UpdateEventAdmin(id, event)
	if err != nil {
		components.ErrorMessage(err.Error()).Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/admin/events")
	w.WriteHeader(http.StatusOK)
}
