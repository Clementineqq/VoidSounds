package service

import (
	"time"
	"voidsounds/internal/domain"
)

type EventService struct{}

func NewEventService() *EventService {
	return &EventService{}
}

// GetAllEvents — пока возвращает мок-данные
func (s *EventService) GetAllEvents() domain.Events {
	return domain.Events{
		{
			ID:          1,
			Title:       "Шум и Выходки в баре «Подвал»",
			Description: "Сольный концерт группы Шум и Выходки. Инди-рок с мощным саундом и неожиданными каверами.",
			Date:        time.Date(2026, 5, 15, 20, 0, 0, 0, time.Local),
			Location:    "Бар «Подвал», МЫТИЩИЫЫЫЫ",
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
	}
}
