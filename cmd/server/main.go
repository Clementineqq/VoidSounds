package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"voidsounds/internal/config"
	"voidsounds/internal/handler"
	"voidsounds/internal/repository"
	"voidsounds/internal/service"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg := config.Load()

	// 2. Подключаемся к базе данных
	repository.InitDB(cfg)

	// 3. Инициализируем репозитории
	eventRepo := repository.NewEventRepository()

	// 4. Инициализируем сервисы
	eventService := service.NewEventService(eventRepo)

	// 5. Инициализируем хендлеры
	eventHandler := handler.NewEventHandler(eventService)

	// 6. Настраиваем роутер
	r := chi.NewRouter()

	// Middleware (промежуточные обработчики)
	r.Use(middleware.Logger)    // Логирование запросов
	r.Use(middleware.Recoverer) // Восстановление после паники
	r.Use(middleware.RequestID) // Добавляет ID каждому запросу

	// 7. Регистрируем маршруты
	r.Get("/", eventHandler.Home)
	r.Get("/events", eventHandler.GetAllEvents)
	r.Get("/event/{id}", eventHandler.GetEventByID)

	// 8. Запускаем сервер
	log.Printf("VoidSounds запущен на http://localhost:%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
