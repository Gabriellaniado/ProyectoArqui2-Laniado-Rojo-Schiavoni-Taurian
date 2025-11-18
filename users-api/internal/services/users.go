package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"users-api/internal/domain"
	"users-api/internal/utils"
)

// UsersService define la lógica de negocio para Users
type UsersRepository interface {
	List(ctx context.Context) ([]domain.UserResponse, error)
	Create(ctx context.Context, user domain.User) (domain.UserResponse, error)
	GetByID(ctx context.Context, id int) (domain.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, id int, user domain.User) (domain.UserResponse, error)
	Delete(ctx context.Context, id int) error
	//Login(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error)
}

// UsersServiceImpl implementa UsersService
// UsersServiceImpl implementa UsersService
type UsersServiceImpl struct {
	repository UsersRepository
}

// Definiciones de errores especificos
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUserID      = errors.New("invalid user ID format")
	ErrEmailRequired      = errors.New("email is required and cannot be empty")
	ErrPasswordRequired   = errors.New("password is required for new users")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFirstNameRequired  = errors.New("first name is required and cannot be empty")
	ErrLastNameRequired   = errors.New("last name is required and cannot be empty")
)

// NewUsersService crea una nueva instancia del service
func NewUsersService(repository UsersRepository) *UsersServiceImpl {
	return &UsersServiceImpl{
		repository: repository,
	}
}

// List obtiene todos los usuarios
func (s *UsersServiceImpl) List(ctx context.Context) ([]domain.UserResponse, error) {
	users, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil

}

// Create valida y crea un nuevo usuario
func (s *UsersServiceImpl) Create(ctx context.Context, user domain.User) (domain.UserResponse, error) {
	// 1. Llamar al validador (Reglas de negocio primarias)
	if err := s.validateCreateUser(user); err != nil {
		return domain.UserResponse{}, err
	}

	// 2. Hash de la contraseña antes de guardar
	user.Password = utils.HashSHA256(user.Password)

	// 3. Intento de Persistencia
	response, err := s.repository.Create(ctx, user)
	// 4. Manejo y traducción de errores del Repositorio (la clave)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "email") {
			return domain.UserResponse{}, ErrEmailAlreadyExists // 409
		}
		if strings.Contains(err.Error(), "email is required") {
			return domain.UserResponse{}, ErrEmailRequired // 400
		}
		if strings.Contains(err.Error(), "first name is required") {
			return domain.UserResponse{}, ErrFirstNameRequired // 400
		}
		if strings.Contains(err.Error(), "last name is required") {
			return domain.UserResponse{}, ErrLastNameRequired // 400
		}

		// Devolvemos cualquier otro error (ej. fallo de conexión a DB)
		return domain.UserResponse{}, err
	}

	// 5. Manejo error
	return response, nil
}

// GetByID obtiene un usuario por su ID
func (s *UsersServiceImpl) GetByID(ctx context.Context, id string) (domain.UserResponse, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return domain.UserResponse{}, ErrInvalidUserID
	}

	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return domain.UserResponse{}, ErrUserNotFound // 404
		}
		return domain.UserResponse{}, err // 500 para otros errores
	}

	return user, nil
}

// GetByEmail obtiene un usuario por su email
func (s *UsersServiceImpl) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if strings.TrimSpace(email) == "" {
		return domain.User{}, ErrEmailRequired
	}
	return s.repository.GetByEmail(ctx, email)
}

// Update actualiza un usuario existente
func (s *UsersServiceImpl) Update(ctx context.Context, id string, user domain.User) (domain.UserResponse, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return domain.UserResponse{}, ErrInvalidUserID //400
	}

	if err := s.validateUser(user); err != nil {
		return domain.UserResponse{}, err
	}

	// Si se envía una nueva contraseña, hashearla
	if user.Password != "" {
		user.Password = utils.HashSHA256(user.Password)
	}

	updatedUser, err := s.repository.Update(ctx, userID, user)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return domain.UserResponse{}, ErrUserNotFound // 404
		}
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "email") {
			return domain.UserResponse{}, ErrEmailAlreadyExists // 409
		}
		return domain.UserResponse{}, err // 500
	}

	return updatedUser, nil
}

// Delete elimina un usuario por ID
func (s *UsersServiceImpl) Delete(ctx context.Context, id string) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return ErrInvalidUserID //400
	}

	err = s.repository.Delete(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ErrUserNotFound // 404
		}
		return err // 500
	}

	return nil
}

// Login valida credenciales de usuario
func (s *UsersServiceImpl) Login(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error) {
	if strings.TrimSpace(loginReq.Email) == "" {
		return domain.LoginResponse{}, ErrEmailRequired //400
	}
	if strings.TrimSpace(loginReq.Password) == "" {
		return domain.LoginResponse{}, ErrPasswordRequired //400
	}

	// Buscar usuario por email
	userModel, err := s.repository.GetByEmail(ctx, loginReq.Email)
	if err != nil {
		return domain.LoginResponse{}, ErrInvalidCredentials //401
	}

	// Verificar contraseña
	if utils.HashSHA256(loginReq.Password) != userModel.Password {
		return domain.LoginResponse{}, ErrInvalidCredentials //401
	}

	// Generar JWT token
	token, err := utils.GenerateJWT(userModel.ID, userModel.IsAdmin)
	if err != nil {
		return domain.LoginResponse{}, errors.New("failed to generate token") //500
	}

	// Obtener el carrito desde products-api
	cart := s.getCartFromProductsAPI(ctx, userModel.ID, token)

	return domain.LoginResponse{
		Token:      token,
		Name:       userModel.FirstName,
		Surname:    userModel.LastName,
		CustomerID: userModel.ID,
		IsAdmin:    userModel.IsAdmin,
		Cart:       cart,
	}, nil
}

// validateUser aplica reglas de negocio para validar un usuario
func (s *UsersServiceImpl) validateUser(user domain.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return ErrEmailRequired // 400
	}
	if strings.TrimSpace(user.FirstName) == "" {
		return ErrFirstNameRequired // 400
	}
	if strings.TrimSpace(user.LastName) == "" {
		return ErrLastNameRequired // 400
	}
	return nil
}

// En UsersServiceImpl o un validador auxiliar en la capa de Servicio/Dominio
func (s *UsersServiceImpl) validateCreateUser(user domain.User) error {
	// ⚠️ Idealmente, el Controller (binding) atrapa esto, pero lo mantenemos por seguridad.
	if err := s.validateUser(user); err != nil {
		return err
	}
	if strings.TrimSpace(user.Password) == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (s *UsersServiceImpl) VerifyToken(token string) error {
	err := utils.ValidateJWT(token)
	if err != nil {
		log.Println("Error al verificar el token")
		return fmt.Errorf("failed to verify token: %w", err)
	}
	return nil
}

func (s *UsersServiceImpl) VerifyAdminToken(token string) error {
	err := utils.ValidateAdminJWT(token)
	if err != nil {
		log.Println("Error al verificar el token de admin")
		return fmt.Errorf("failed to verify admin token: %w", err)
	}
	return nil
}

// getCartFromProductsAPI obtiene el carrito del usuario desde products-api
func (s *UsersServiceImpl) getCartFromProductsAPI(ctx context.Context, customerID int, token string) interface{} {
	// URL del products-api (ajustar según tu configuración de docker-compose)
	url := fmt.Sprintf("http://products-api:8081/carrito/%d", customerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating cart request: %v", err)
		return map[string]interface{}{"items": []interface{}{}, "total": 0, "item_count": 0}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling products-api for cart: %v", err)
		return map[string]interface{}{"items": []interface{}{}, "total": 0, "item_count": 0}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Products-api returned status %d for cart", resp.StatusCode)
		return map[string]interface{}{"items": []interface{}{}, "total": 0, "item_count": 0}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading cart response: %v", err)
		return map[string]interface{}{"items": []interface{}{}, "total": 0, "item_count": 0}
	}

	var cart interface{}
	if err := json.Unmarshal(body, &cart); err != nil {
		log.Printf("Error parsing cart response: %v", err)
		return map[string]interface{}{"items": []interface{}{}, "total": 0, "item_count": 0}
	}

	return cart
}
