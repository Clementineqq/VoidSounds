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

func (s *UserService) Register(req *domain.RegisterRequest) (*domain.User, error) {

	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("некорректный email")
	}

	existing, _ := s.repo.GetByEmail(req.Email)
	if existing != nil {
		return nil, fmt.Errorf("пользователь с таким email уже существует")
	}

	if len(req.Password) < 6 {
		return nil, fmt.Errorf("пароль должен быть не менее 6 символов")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "viewer"
	}

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

	return user, nil
}

func (s *UserService) Login(email, password string) (*domain.User, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	return user, nil
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func isValidEmail(email string) bool {

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) ChangeUserRole(userID int, role string) error {
	return s.repo.ChangeUserRole(userID, role)
}

func (s *UserService) BanUser(userID int) error {
	return s.repo.BanUser(userID)
}
