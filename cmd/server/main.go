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
	mymw.InitSessionStore("super-secret-key-change-me")

	// 4. Инициализируем репозитории
	eventRepo := repository.NewEventRepository()
	userRepo := repository.NewUserRepository()

	// 5. Инициализируем сервисы
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)

	// 6. Инициализируем хендлеры
	eventHandler := handler.NewEventHandler(eventService)
	authHandler := handler.NewAuthHandler(userService)

	// 7. Настраиваем роутер
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mymw.AuthMiddleware) // Добавляем middleware авторизации

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

	// Защищенные маршруты (только для авторизованных)
	r.Group(func(r chi.Router) {
		// TODO: добавить middleware проверки авторизации
		// r.Use(RequireAuth)

		// Покупка билетов
		// r.Post("/event/{id}/buy", ticketHandler.BuyTicket)

		// Личный кабинет
		// r.Get("/profile", profileHandler.Profile)
	})

	// Маршруты организатора
	r.Group(func(r chi.Router) {
		// TODO: добавить middleware проверки роли organizer
		// r.Use(RequireRole("organizer"))

		// Создание мероприятия
		// r.Get("/organizer/events/create", eventHandler.ShowCreateForm)
		// r.Post("/organizer/events", eventHandler.CreateEvent)
		// r.Get("/organizer/events", eventHandler.GetOrganizerEvents)
		// r.Put("/organizer/events/{id}", eventHandler.UpdateEvent)
		// r.Delete("/organizer/events/{id}", eventHandler.DeleteEvent)
	})

	// 9. Запускаем сервер
	log.Printf("VoidSounds запущен на http://localhost:%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
