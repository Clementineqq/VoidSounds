package handler

import (
	"net/http"
	"voidsounds/internal/components"
	"voidsounds/internal/middleware"
	"voidsounds/internal/service"
)

type AdminHandler struct {
	eventService *service.EventService
	userService  *service.UserService
}

func NewAdminHandler(eventService *service.EventService, userService *service.UserService) *AdminHandler {
	return &AdminHandler{eventService: eventService, userService: userService}
}

// GET /admin - панель администратора
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == 0 {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Простая статистика (можно расширить)
	stats := map[string]int{
		"users":   42,  // заглушка
		"events":  15,  // заглушка
		"tickets": 128, // заглушка
	}

	components.AdminDashboard(stats).Render(r.Context(), w)
}
