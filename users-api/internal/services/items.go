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
// Capa intermedia entre Controllers (HTTP) y Repository (datos)
// Responsabilidades: validaciones, transformaciones, reglas de negocio
type UsersService interface {
	// List retorna todos los usuarios
	List(ctx context.Context) ([]domain.UserDto, error)

	// Create valida y crea un nuevo usuario
	Create(ctx context.Context, user domain.UserDto) (domain.UserDto, error)

	// GetByID obtiene un usuario por su ID
	GetByID(ctx context.Context, id string) (domain.UserDto, error)

	// GetByEmail obtiene un usuario por su email
	GetByEmail(ctx context.Context, email string) (domain.UserDto, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, id string, user domain.UserDto) (domain.UserDto, error)

	// Delete elimina un usuario por ID
	Delete(ctx context.Context, id string) error

	// Login valida credenciales de usuario
	Login(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error)
}

// UsersServiceImpl implementa UsersService
type UsersServiceImpl struct {
	repository repository.UsersRepository // Inyección de dependencia
}

// NewUsersService crea una nueva instancia del service
func NewUsersService(repository repository.UsersRepository) UsersService {
	return &UsersServiceImpl{repository: repository}
}

// List obtiene todos los usuarios
func (s *UsersServiceImpl) List(ctx context.Context) ([]domain.UserDto, error) {
	return s.repository.List(ctx)
}

// Create valida y crea un nuevo usuario
func (s *UsersServiceImpl) Create(ctx context.Context, user domain.UserDto) (domain.UserDto, error) {
	if err := s.validateUser(user); err != nil {
		return domain.UserDto{}, err
	}

	// Hash de la contraseña antes de guardar
	user.Password = utils.HashSHA256(user.Password)

	return s.repository.Create(ctx, user)
}

// GetByID obtiene un usuario por su ID
func (s *UsersServiceImpl) GetByID(ctx context.Context, id string) (domain.UserDto, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return domain.UserDto{}, errors.New("invalid user ID format")
	}
	return s.repository.GetByID(ctx, userID)
}

// GetByEmail obtiene un usuario por su email
func (s *UsersServiceImpl) GetByEmail(ctx context.Context, email string) (domain.UserDto, error) {
	if strings.TrimSpace(email) == "" {
		return domain.UserDto{}, errors.New("email is required")
	}
	return s.repository.GetByEmail(ctx, email)
}

// Update actualiza un usuario existente
func (s *UsersServiceImpl) Update(ctx context.Context, id string, user domain.UserDto) (domain.UserDto, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return domain.UserDto{}, errors.New("invalid user ID format")
	}

	if err := s.validateUser(user); err != nil {
		return domain.UserDto{}, err
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
	user, err := s.repository.GetByEmail(ctx, loginReq.Email)
	if err != nil {
		return domain.LoginResponse{}, errors.New("invalid credentials")
	}

	// Verificar contraseña
	hashedPassword := utils.HashSHA256(loginReq.Password)
	if user.Password != hashedPassword {
		return domain.LoginResponse{}, errors.New("invalid credentials")
	}

	// TODO: Generar JWT token aquí
	token, err := utils.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		return domain.LoginResponse{}, errors.New("failed to generate token")
	}

	return domain.LoginResponse{
		Token:   token,
		Name:    user.FirstName,
		Surname: user.LastName,
	}, nil
}

// validateUser aplica reglas de negocio para validar un usuario
func (s *UsersServiceImpl) validateUser(user domain.UserDto) error {
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required and cannot be empty")
	}
	if strings.TrimSpace(user.FirstName) == "" {
		return errors.New("first name is required and cannot be empty")
	}
	if strings.TrimSpace(user.LastName) == "" {
		return errors.New("last name is required and cannot be empty")
	}
	if user.ID == 0 && strings.TrimSpace(user.Password) == "" {
		return errors.New("password is required for new users")
	}
	return nil
}
