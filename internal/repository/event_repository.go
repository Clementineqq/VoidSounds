package repository

import (
	"time"
	"voidsounds/internal/domain"
)

type EventRepository interface {
	GetAll() (domain.Events, error)
	GetByID(id int) (*domain.Event, error)
	// позже добавим Create, Update, Delete
}

type eventRepository struct {
	// пока пусто, потом будет подключение к БД
}

func NewEventRepository() EventRepository {
	return &eventRepository{}
}

func (r *eventRepository) GetAll() (domain.Events, error) {
	// Пока возвращаем мок-данные (позже заменим на запрос к БД)
	return domain.Events{
		{
			ID:          1,
			Title:       "Шум и Выходки в баре «Подвал»",
			Description: "Сольный концерт группы Шум и Выходки. Инди-рок с мощным саундом.",
			Date:        time.Date(2026, 5, 15, 20, 0, 0, 0, time.Local),
			Location:    "Бар «Подвал», Москва",
			Genre:       "Инди-рок",
			Price:       1500,
			Available:   87,
		},
		{
			ID:          2,
			Title:       "Электронная ночь на крыше",
			Description: "Три артиста электронной сцены. Живая электроника и визуальное шоу.",
			Date:        time.Date(2026, 5, 22, 22, 0, 0, 0, time.Local),
			Location:    "Крыша «Flora», Санкт-Петербург",
			Genre:       "Электронная",
			Price:       2000,
			Available:   45,
		},
	}, nil
}

func (r *eventRepository) GetByID(id int) (*domain.Event, error) {
	// Пока заглушка
	return nil, nil
}
