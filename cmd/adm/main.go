package main

import (
	"fmt"
	"log"

	"voidsounds/internal/config"
	"voidsounds/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	repository.InitDB(cfg)

	email := "admin@voidsounds.ru"
	password := "admin"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	query := `
		INSERT INTO users (email, password_hash, name, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE 
		SET role = 'admin', password_hash = $2, name = $3
	`

	_, err = repository.DB.Exec(query, email, string(hashedPassword), "Админ", "admin")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" Email: %s\n", email)
	fmt.Printf("  Пароль: %s\n", password)
}
