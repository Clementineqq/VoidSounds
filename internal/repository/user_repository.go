package repository

import (
	"fmt"
	"voidsounds/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id int) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)          // ← ДОБАВИТЬ
	ChangeUserRole(userID int, role string) error // ← ДОБАВИТЬ
	BanUser(userID int) error                     // ← ДОБАВИТЬ
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *domain.User) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}

	query := `
        INSERT INTO users (email, password_hash, name, role)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `

	err := DB.QueryRowx(
		query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не подключена")
	}

	query := `
        SELECT id, email, password_hash, name, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user domain.User
	err := DB.Get(&user, query, email)
	if err != nil {
		return nil, fmt.Errorf("пользователь с email %s не найден: %w", email, err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не подключена")
	}

	query := `
        SELECT id, email, name, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user domain.User
	err := DB.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("пользователь с ID %d не найден: %w", id, err)
	}

	return &user, nil
}

// GetAllUsers - получить всех пользователей
func (r *userRepository) GetAllUsers() ([]domain.User, error) {
	if DB == nil {
		return []domain.User{}, nil
	}
	query := `SELECT id, email, name, role, created_at FROM users ORDER BY created_at DESC`
	var users []domain.User
	err := DB.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователей: %w", err)
	}
	return users, nil
}

// ChangeUserRole - сменить роль пользователя
func (r *userRepository) ChangeUserRole(userID int, role string) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}
	query := `UPDATE users SET role = $1 WHERE id = $2`
	result, err := DB.Exec(query, role, userID)
	if err != nil {
		return fmt.Errorf("ошибка смены роли: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}

func (r *userRepository) BanUser(userID int) error {
	if DB == nil {
		return fmt.Errorf("база данных не подключена")
	}

	var isBanned bool
	err := DB.Get(&isBanned, `SELECT is_banned FROM users WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	query := `UPDATE users SET is_banned = $1 WHERE id = $2`
	_, err = DB.Exec(query, !isBanned, userID)
	if err != nil {
		return fmt.Errorf("ошибка бана: %w", err)
	}
	return nil
}
