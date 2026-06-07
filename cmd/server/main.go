package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"voidsounds/internal/config"
	"voidsounds/internal/handler"
	mymw "voidsounds/internal/middleware"
	"voidsounds/internal/repository"
	"voidsounds/internal/service"
)

func main() {
	cfg := config.Load()

	repository.InitDB(cfg)

	mymw.InitSessionStore(cfg.SessionSecret)

	eventRepo := repository.NewEventRepository()
	userRepo := repository.NewUserRepository()

	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)

	eventHandler := handler.NewEventHandler(eventService, userService)
	authHandler := handler.NewAuthHandler(userService)
	adminHandler := handler.NewAdminHandler(eventService, userService)
	pageHandler := handler.NewPageHandler()

	r := chi.NewRouter()

	r.Use(mymw.MethodOverride)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mymw.AuthMiddleware)

	r.Get("/", eventHandler.Home)
	r.Get("/events", eventHandler.GetAllEvents)
	r.Get("/event/{id}", eventHandler.GetEventByID)

	r.Get("/register", authHandler.ShowRegister)
	r.Post("/register", authHandler.Register)
	r.Get("/login", authHandler.ShowLogin)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Get("/auth/status", authHandler.GetAuthStatus)

	r.Group(func(r chi.Router) {
		r.Use(mymw.RequireAuth)
		r.Post("/event/{id}/buy", eventHandler.BuyTicket)
		r.Get("/profile", eventHandler.Profile)
		r.Get("/ticket/{id}/qr", eventHandler.TicketQR)
	})

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

	r.Group(func(r chi.Router) {
		r.Use(mymw.RequireAuth, mymw.RequireRole("admin"))
		r.Get("/admin", adminHandler.Dashboard)
		r.Get("/admin/users", adminHandler.GetUsers)
		r.Get("/admin/users/{id}", adminHandler.GetUserByID)
		r.Post("/admin/users/{id}/role", adminHandler.ChangeUserRole)
		r.Post("/admin/users/{id}/ban", adminHandler.BanUser)
		r.Get("/admin/events", adminHandler.GetEvents)
		r.Get("/admin/events/{id}/edit", adminHandler.EditEventForm)
		r.Post("/admin/events/{id}/update", adminHandler.UpdateEvent)
		r.Post("/admin/events/{id}/status", adminHandler.ChangeEventStatus)
		r.Delete("/admin/events/{id}", adminHandler.DeleteEvent)
	})

	r.Get("/for-organizers", pageHandler.ForOrganizers)
	r.Get("/organizer/{id}", eventHandler.ShowOrganizerProfile)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	log.Printf("VoidSounds запущен на http://localhost:%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
