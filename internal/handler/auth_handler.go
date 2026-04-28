package handler

import (
	"net/http"

	"voidsounds/internal/components"
	"voidsounds/internal/domain"
	"voidsounds/internal/middleware"
	"voidsounds/internal/service"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// GET /register - показать форму регистрации
func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	component := components.RegisterForm()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /register - обработка регистрации
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Парсим форму
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	req := &domain.RegisterRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Role:     r.FormValue("role"),
	}

	// Регистрируем пользователя
	user, err := h.userService.Register(req)
	if err != nil {
		// Показываем ошибку
		component := components.ErrorMessage(err.Error())
		component.Render(r.Context(), w)
		return
	}

	// Создаем сессию
	session, _ := middleware.Store.Get(r, "user-session")
	session.Values["user_id"] = user.ID
	session.Values["user_email"] = user.Email
	session.Values["user_name"] = user.Name
	session.Values["user_role"] = user.Role
	session.Save(r, w)

	// Показываем успех
	component := components.AuthSuccess("Регистрация прошла успешно!", "/events")
	component.Render(r.Context(), w)
}

// GET /login - показать форму входа
func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	component := components.LoginForm()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /login - обработка входа
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Парсим форму
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Проверяем учетные данные
	user, err := h.userService.Login(email, password)
	if err != nil {
		component := components.ErrorMessage("Неверный email или пароль")
		component.Render(r.Context(), w)
		return
	}

	// Создаем сессию
	session, _ := middleware.Store.Get(r, "user-session")
	session.Values["user_id"] = user.ID
	session.Values["user_email"] = user.Email
	session.Values["user_name"] = user.Name
	session.Values["user_role"] = user.Role
	session.Save(r, w)

	// Показываем успех
	component := components.AuthSuccess("Добро пожаловать!", "/events")
	component.Render(r.Context(), w)
}

// GET /logout - выход из системы
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Удаляем сессию
	session, _ := middleware.Store.Get(r, "user-session")
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Перенаправляем на главную
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

// GET /auth/status - получить текущий статус авторизации
func (h *AuthHandler) GetAuthStatus(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "user-session")

	userID, ok := session.Values["user_id"].(int)
	if !ok || userID == 0 {
		// Не авторизован
		component := components.AuthStatus(false, "", "")
		component.Render(r.Context(), w)
		return
	}

	// Авторизован
	userName, _ := session.Values["user_name"].(string)
	userRole, _ := session.Values["user_role"].(string)

	component := components.AuthStatus(true, userName, userRole)
	component.Render(r.Context(), w)
}
