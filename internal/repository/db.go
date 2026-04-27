package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"voidsounds/internal/config"
)

var DB *sqlx.DB

func InitDB(cfg *config.Config) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Printf("Ошибка подключения к PostgreSQL: %v", err)
		log.Println("Проверьте, что контейнер PostgreSQL запущен (docker compose up -d postgres)")
		DB = nil
		return
	}

	// Проверка соединения
	if err = DB.Ping(); err != nil {
		log.Printf("Ошибка Ping к PostgreSQL: %v", err)
		DB = nil
		return
	}

	log.Println("✅ Успешно подключились к PostgreSQL")
}
