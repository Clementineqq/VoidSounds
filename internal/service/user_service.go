package service

import (
	"fmt"
	"regexp"
	"strings"

	"voidsounds/internal/domain"
	"voidsounds/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Регистрация пользователя
func (s *UserService) Register(req *domain.RegisterRequest) (*domain.User, error) {
	// Валидация email
	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("некорректный email")
	}

	// Проверка, что email не занят
	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, fmt.Errorf("пользователь с таким email уже существует")
	}

	// Валидация пароля (минимум 6 символов)
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("пароль должен быть не менее 6 символов")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	// Если роль не указана, ставим "viewer"
	role := req.Role
	if role == "" {
		role = "viewer"
	}

	// Только админ может создать организатора (пока разрешим всем, потом ограничим)
	// Для MVP можно разрешить любую роль, но в продакшене нужно проверять

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Role:         role,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Возвращаем пользователя без пароля
	return user, nil
}

// Вход пользователя
func (s *UserService) Login(email, password string) (*domain.User, error) {
	// Находим пользователя по email
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	return user, nil
}

// Получить пользователя по ID
func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.repo.GetByID(id)
}

// Вспомогательная функция проверки email
func isValidEmail(email string) bool {
	// Простая проверка email через regex
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}
