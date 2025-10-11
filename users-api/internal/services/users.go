package services

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"users-api/internal/domain"
	"users-api/internal/repository"
	"users-api/internal/utils"
)

// UsersService define la lógica de negocio para Users
type UsersService interface {
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error) // Cambiar a string
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, id string, user domain.User) (domain.User, error) // Cambiar a string
	Delete(ctx context.Context, id string) error                                  // Cambiar a string
	Login(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error)
}

// UsersServiceImpl implementa UsersService
// UsersServiceImpl implementa UsersService
type UsersServiceImpl struct {
	repository repository.UsersRepository
}

// NewUsersService crea una nueva instancia del service
func NewUsersService(repository repository.UsersRepository) *UsersServiceImpl {
	return &UsersServiceImpl{
		repository: repository,
	}
}

// List obtiene todos los usuarios
func (s *UsersServiceImpl) List(ctx context.Context) ([]domain.User, error) {
	return s.repository.List(ctx)
}

// Create valida y crea un nuevo usuario
func (s *UsersServiceImpl) Create(ctx context.Context, user domain.User) (domain.User, error) {
	if err := s.validateCreateUser(user); err != nil {
		return domain.User{}, err
	}

	// Hash de la contraseña antes de guardar
	user.Password = utils.HashSHA256(user.Password)

	return s.repository.Create(ctx, user)
}

// GetByID obtiene un usuario por su ID
func (s *UsersServiceImpl) GetByID(ctx context.Context, id string) (domain.User, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return domain.User{}, errors.New("invalid user ID format")
	}
	return s.repository.GetByID(ctx, userID)
}

// GetByEmail obtiene un usuario por su email
func (s *UsersServiceImpl) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if strings.TrimSpace(email) == "" {
		return domain.User{}, errors.New("email is required")
	}
	return s.repository.GetByEmail(ctx, email)
}

// Update actualiza un usuario existente
func (s *UsersServiceImpl) Update(ctx context.Context, id string, user domain.User) (domain.User, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return domain.User{}, errors.New("invalid user ID format")
	}

	if err := s.validateUser(user); err != nil {
		return domain.User{}, err
	}

	// Si se envía una nueva contraseña, hashearla
	if user.Password != "" {
		user.Password = utils.HashSHA256(user.Password)
	}

	return s.repository.Update(ctx, userID, user)
}

// Delete elimina un usuario por ID
func (s *UsersServiceImpl) Delete(ctx context.Context, id string) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	return s.repository.Delete(ctx, userID)
}

// Login valida credenciales de usuario
func (s *UsersServiceImpl) Login(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error) {
	if strings.TrimSpace(loginReq.Email) == "" {
		return domain.LoginResponse{}, errors.New("email is required")
	}
	if strings.TrimSpace(loginReq.Password) == "" {
		return domain.LoginResponse{}, errors.New("password is required")
	}

	// Buscar usuario por email
	userModel, err := s.repository.GetByEmail(ctx, loginReq.Email)
	if err != nil {
		return domain.LoginResponse{}, errors.New("invalid credentials")
	}

	// Verificar contraseña
	if utils.HashSHA256(loginReq.Password) != userModel.Password {
		return domain.LoginResponse{}, errors.New("invalid credentials")
	}

	// Generar JWT token
	token, err := utils.GenerateJWT(userModel.ID, userModel.IsAdmin)
	if err != nil {
		return domain.LoginResponse{}, errors.New("failed to generate token")
	}

	return domain.LoginResponse{
		Token:   token,
		Name:    userModel.FirstName,
		Surname: userModel.LastName,
	}, nil
}

// validateUser aplica reglas de negocio para validar un usuario
func (s *UsersServiceImpl) validateUser(user domain.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required and cannot be empty")
	}
	if strings.TrimSpace(user.FirstName) == "" {
		return errors.New("first name is required and cannot be empty")
	}
	if strings.TrimSpace(user.LastName) == "" {
		return errors.New("last name is required and cannot be empty")
	}
	return nil
}

// Crear método separado para validar CREATE
func (s *UsersServiceImpl) validateCreateUser(user domain.User) error {
	if err := s.validateUser(user); err != nil {
		return err
	}
	// Solo requerir password en CREATE
	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is required for new users")
	}
	return nil
}
