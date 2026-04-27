package main

import (
	"log"
	"os"

	"voidsounds/internal/config"
	"voidsounds/internal/repository"
)

func main() {
	cfg := config.Load()
	repository.InitDB(cfg)

	if repository.DB == nil {
		log.Fatal("Не удалось подключиться к базе данных")
	}

	migrationFiles := []string{
		"migrations/002_create_full_schema.sql",
		"migrations/003_seed_initial_data.sql",
	}

	for _, file := range migrationFiles {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Не удалось прочитать файл %s: %v", file, err)
			continue
		}

		_, err = repository.DB.Exec(string(sqlBytes))
		if err != nil {
			log.Printf("Ошибка выполнения миграции %s: %v", file, err)
			continue
		}

		log.Printf("Миграция %s успешно применена", file)
	}

	log.Println(" Все миграции успешно выполнены, пьем пиво!")
}
