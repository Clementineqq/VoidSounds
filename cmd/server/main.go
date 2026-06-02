package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"voidsounds/internal/config"
	"voidsounds/internal/handler"
	mymw "voidsounds/internal/middleware" // alias чтобы не конфликтовало с chi/middleware
	"voidsounds/internal/repository"
	"voidsounds/internal/service"
)

func main() {

	// 1. Загружаем конфигурацию
	cfg := config.Load()

	// 2. Подключаемся к базе данных
	repository.InitDB(cfg)

	// 3. Инициализируем сессии
	mymw.InitSessionStore(cfg.SessionSecret)
	// 4. Инициализируем репозитории
	eventRepo := repository.NewEventRepository()
	userRepo := repository.NewUserRepository()

	// 5. Инициализируем сервисы
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)

	// 6. Инициализируем хендлеры
	eventHandler := handler.NewEventHandler(eventService, userService)
	authHandler := handler.NewAuthHandler(userService)
	adminHandler := handler.NewAdminHandler(eventService, userService)
	// 7. Настраиваем роутер
	r := chi.NewRouter()

	// Middleware
	r.Use(mymw.MethodOverride)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mymw.AuthMiddleware)

	// 8. Регистрируем маршруты

	// Публичные маршруты
	r.Get("/", eventHandler.Home)
	r.Get("/events", eventHandler.GetAllEvents)
	r.Get("/event/{id}", eventHandler.GetEventByID)

	// Маршруты авторизации
	r.Get("/register", authHandler.ShowRegister)
	r.Post("/register", authHandler.Register)
	r.Get("/login", authHandler.ShowLogin)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Get("/auth/status", authHandler.GetAuthStatus)

	// Защищённые маршруты (только для авторизованных)
	r.Group(func(r chi.Router) {
		r.Use(mymw.RequireAuth)
		r.Post("/event/{id}/buy", eventHandler.BuyTicket)
		r.Get("/profile", eventHandler.Profile)
		r.Get("/ticket/{id}/qr", eventHandler.TicketQR)
	})

	// Маршруты организатора (авторизация + роль organizer)
	r.Group(func(r chi.Router) {
		r.Use(mymw.RequireAuth, mymw.RequireRole("organizer"))
		r.Post("/organizer/events/{id}/status", eventHandler.ChangeStatus)
		r.Get("/organizer/events/create", eventHandler.ShowCreateForm)
		r.Post("/organizer/events", eventHandler.CreateEvent)
		r.Get("/organizer/events", eventHandler.GetOrganizerEvents)
		r.Delete("/organizer/events/{id}", eventHandler.DeleteEvent)
		r.Get("/organizer/events/{id}/edit", eventHandler.ShowEditForm)
		r.Post("/organizer/events/{id}/update", eventHandler.UpdateEvent)
	})
	// Админка (только для роли admin)
	r.Group(func(r chi.Router) {
		r.Use(mymw.RequireAuth, mymw.RequireRole("admin"))
		r.Get("/admin", adminHandler.Dashboard)
	})

	// Инициализация хендлера страниц
	pageHandler := handler.NewPageHandler()

	r.Get("/organizer/{id}", eventHandler.ShowOrganizerProfile)

	// Публичные информационные страницы
	r.Get("/artists", pageHandler.Artists)
	r.Get("/for-organizers", pageHandler.ForOrganizers)

	log.Printf("VoidSounds запущен на http://localhost:%s", cfg.ServerPort)
	// Раздача загруженных файлов
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
