package main

import (
	"log"
	"os"

	"voidsounds/internal/config"
	"voidsounds/internal/repository"
)

func main() {
	cfg := config.Load()

	// Принудительно пытаемся подключиться
	repository.InitDB(cfg)

	if repository.DB == nil {
		log.Fatal(" Не удалось подключиться к базе данных. Убедитесь, что PostgreSQL запущен через Docker.")
	}

	// Выполняем миграцию
	sqlBytes, err := os.ReadFile("migrations/001_create_events_table.sql")
	if err != nil {
		log.Fatal("Не удалось прочитать файл миграции:", err)
	}

	_, err = repository.DB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Ошибка выполнения миграции:", err)
	}

	log.Println("  Миграция успешно применена! Таблица events создана.")
}
