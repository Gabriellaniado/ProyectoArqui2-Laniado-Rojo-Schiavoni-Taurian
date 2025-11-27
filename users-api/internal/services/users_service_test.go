package services

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"users-api/internal/domain"

	"github.com/h2non/gock"
)

// ============================================
// MOCKS version falsa que simula el repositorio y la bd
// ============================================

// MockUsersRepository simula el repositorio de usuarios
type MockUsersRepository struct {
	users      map[int]domain.User // Almacena usuarios en memoria (no en MySQL)
	nextID     int                 // Simula auto-increment de la BD
	shouldFail bool                // Para simular errores
}

func NewMockUsersRepository() *MockUsersRepository {
	return &MockUsersRepository{
		users:  make(map[int]domain.User),
		nextID: 1,
	}
}

// Simula la obtención de todos los usuarios y respuesta de arreglo de UserResponse
func (m *MockUsersRepository) List(ctx context.Context) ([]domain.UserResponse, error) {
	if m.shouldFail {
		return nil, errors.New("database error")
	}

	var result []domain.UserResponse
	for _, user := range m.users {
		result = append(result, domain.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsAdmin:   user.IsAdmin,
		})
	}
	return result, nil
}

// Simula la creación de un usuario y respuesta de UserResponse
func (m *MockUsersRepository) Create(ctx context.Context, user domain.User) (domain.UserResponse, error) {
	// permite simular un error de base de datos
	if m.shouldFail {
		return domain.UserResponse{}, errors.New("database error")
	}

	// Verificar si el email ya existe (constraint UNIQUE)
	for _, u := range m.users {
		if u.Email == user.Email {
			return domain.UserResponse{}, errors.New("Duplicate entry for email")
		}
	}
	// Insertar un usuario en memoria
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user

	return domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	}, nil
}

// Simula la obtención de un usuario por ID y respuesta de UserResponse
func (m *MockUsersRepository) GetByID(ctx context.Context, id int) (domain.UserResponse, error) {
	if m.shouldFail {
		return domain.UserResponse{}, errors.New("database error")
	}

	user, exists := m.users[id]
	if !exists {
		return domain.UserResponse{}, errors.New("user not found")
	}

	return domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	}, nil
}

// Simula la obtención de un usuario por email y respuesta de domain User
func (m *MockUsersRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if m.shouldFail {
		return domain.User{}, errors.New("database error")
	}

	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, errors.New("user not found")
}

// Simula la actualización de un usuario y respuesta de UserResponse
func (m *MockUsersRepository) Update(ctx context.Context, id int, user domain.User) (domain.UserResponse, error) {
	if m.shouldFail {
		return domain.UserResponse{}, errors.New("database error")
	}

	if _, exists := m.users[id]; !exists {
		return domain.UserResponse{}, errors.New("user not found")
	}

	user.ID = id
	m.users[id] = user

	return domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	}, nil
}

// Simula la eliminación de un usuario y respuesta de error si falla
func (m *MockUsersRepository) Delete(ctx context.Context, id int) error {
	if m.shouldFail {
		return errors.New("database error")
	}

	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	return nil
}

// ============================================
// TESTS patron AAA(Arrange, Act, Assert) (Preparar, Ejecutar, Verificar)
// ============================================

// Para Create se prueba: exito, email duplicado, email vacio, password vacio, first name vacio, last name vacio

// TestCreate_Success verifica la creación exitosa de un usuario
func TestCreate_Success(t *testing.T) {
	// Preparar
	mockRepo := NewMockUsersRepository() // Mock en lugar de MySQL
	service := NewUsersService(mockRepo) // Servicio con el mock

	user := domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		IsAdmin:   false,
	}

	// Ejecutar
	created, err := service.Create(context.Background(), user)

	// Verificar
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if created.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if created.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, created.Email)
	}

	if created.FirstName != user.FirstName {
		t.Errorf("Expected first name %s, got %s", user.FirstName, created.FirstName)
	}
}

// TestCreate_DuplicateEmail verifica que no se permita email duplicado
func TestCreate_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user1 := domain.User{
		Email:     "duplicate@example.com",
		Password:  "password123",
		FirstName: "User",
		LastName:  "One",
	}

	user2 := domain.User{
		Email:     "duplicate@example.com",
		Password:  "password456",
		FirstName: "User",
		LastName:  "Two",
	}

	_, err1 := service.Create(context.Background(), user1)
	_, err2 := service.Create(context.Background(), user2)

	if err1 != nil {
		t.Errorf("First creation should succeed, got error: %v", err1)
	}

	if err2 == nil {
		t.Error("Expected error for duplicate email")
	}

	if !errors.Is(err2, ErrEmailAlreadyExists) {
		t.Errorf("Expected ErrEmailAlreadyExists, got %v", err2)
	}
}

// TestCreate_EmptyEmail verifica validación de email vacío
func TestCreate_EmptyEmail(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	_, err := service.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error for empty email")
	}

	if !errors.Is(err, ErrEmailRequired) {
		t.Errorf("Expected ErrEmailRequired, got %v", err)
	}
}

// TestCreate_EmptyPassword verifica validación de password vacío
func TestCreate_EmptyPassword(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "test@example.com",
		Password:  "",
		FirstName: "John",
		LastName:  "Doe",
	}

	_, err := service.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error for empty password")
	}

	if !errors.Is(err, ErrPasswordRequired) {
		t.Errorf("Expected ErrPasswordRequired, got %v", err)
	}
}

// TestCreate_EmptyFirstName verifica validación de first name vacío
func TestCreate_EmptyFirstName(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "",
		LastName:  "Doe",
	}

	_, err := service.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error for empty first name")
	}

	if !errors.Is(err, ErrFirstNameRequired) {
		t.Errorf("Expected ErrFirstNameRequired, got %v", err)
	}
}

// TestCreate_EmptyLastName verifica validación de last name vacío
func TestCreate_EmptyLastName(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "",
	}

	_, err := service.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error for empty last name")
	}

	if !errors.Is(err, ErrLastNameRequired) {
		t.Errorf("Expected ErrLastNameRequired, got %v", err)
	}
}

//  TESTS PARA LOGIN: exito, contraseña incorrecta, email inexistente

// TestLogin_Success verifica el login exitoso
func TestLogin_Success(t *testing.T) {
	defer gock.Off() // Limpiar interceptores después del test

	//  Interceptar llamadas HTTP a products-api
	gock.New("http://products-api:8081").
		Get("/carrito/1"). // El customerID será 1 (primer usuario creado)
		Reply(200).
		JSON(map[string]interface{}{
			"items":      []interface{}{},
			"total":      0,
			"item_count": 0,
		})

	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	// Crear un usuario primero
	user := domain.User{
		Email:     "login@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	_, _ = service.Create(context.Background(), user)

	loginReq := domain.LoginRequest{
		Email:    "login@example.com",
		Password: "password123",
	}

	// Ejecutar login

	response, err := service.Login(context.Background(), loginReq)

	// Verificar response
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}

	if response.Name != user.FirstName {
		t.Errorf("Expected name %s, got %s", user.FirstName, response.Name)
	}

	if response.Surname != user.LastName {
		t.Errorf("Expected surname %s, got %s", user.LastName, response.Surname)
	}
}

// TestLogin_WrongPassword verifica rechazo de contraseña incorrecta
func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "wrongpass@example.com",
		Password:  "correctpassword",
		FirstName: "John",
		LastName:  "Doe",
	}
	_, _ = service.Create(context.Background(), user)

	loginReq := domain.LoginRequest{
		Email:    "wrongpass@example.com",
		Password: "wrongpassword",
	}

	_, err := service.Login(context.Background(), loginReq)

	if err == nil {
		t.Error("Expected error for wrong password")
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

// TestLogin_UserNotFound verifica rechazo de email inexistente
func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	loginReq := domain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := service.Login(context.Background(), loginReq)

	if err == nil {
		t.Error("Expected error for nonexistent user")
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

// TESTS PARA GET BY ID: exito, no encontrado, id inválido

// TestGetByID_Success verifica obtener usuario por ID
func TestGetByID_Success(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "getbyid@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	created, _ := service.Create(context.Background(), user)

	found, err := service.GetByID(context.Background(), strconv.Itoa(created.ID))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, found.ID)
	}

	if found.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, found.Email)
	}
}

// TestGetByID_NotFound verifica manejo de usuario no encontrado
func TestGetByID_NotFound(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	_, err := service.GetByID(context.Background(), "999")

	if err == nil {
		t.Error("Expected error for nonexistent user")
	}

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

// TestGetByID_InvalidID verifica manejo de ID inválido
func TestGetByID_InvalidID(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	_, err := service.GetByID(context.Background(), "invalid")

	if err == nil {
		t.Error("Expected error for invalid ID")
	}

	if !errors.Is(err, ErrInvalidUserID) {
		t.Errorf("Expected ErrInvalidUserID, got %v", err)
	}
}

// TESTS PARA LIST: exito
// TestList_Success verifica obtener todos los usuarios
func TestList_Success(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user1 := domain.User{Email: "user1@example.com", Password: "pass1", FirstName: "User", LastName: "One"}
	user2 := domain.User{Email: "user2@example.com", Password: "pass2", FirstName: "User", LastName: "Two"}
	_, _ = service.Create(context.Background(), user1)
	_, _ = service.Create(context.Background(), user2)

	users, err := service.List(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

// TEST PARA UPDATE: exito
// TestUpdate_Success verifica actualización de usuario
func TestUpdate_Success(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "original@example.com",
		Password:  "password123",
		FirstName: "Original",
		LastName:  "Name",
	}
	created, _ := service.Create(context.Background(), user)

	updateReq := domain.User{
		Email:     "updated@example.com",
		Password:  "newpassword",
		FirstName: "Updated",
		LastName:  "Name",
	}

	updated, err := service.Update(context.Background(), strconv.Itoa(created.ID), updateReq)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if updated.FirstName != "Updated" {
		t.Errorf("Expected first name 'Updated', got '%s'", updated.FirstName)
	}
}

// TEST PARA DELETE: exito
// TestDelete_Success verifica eliminación de usuario
func TestDelete_Success(t *testing.T) {
	mockRepo := NewMockUsersRepository()
	service := NewUsersService(mockRepo)

	user := domain.User{
		Email:     "delete@example.com",
		Password:  "password123",
		FirstName: "Delete",
		LastName:  "User",
	}

	// Crear un usuario primero
	created, _ := service.Create(context.Background(), user)

	// Eliminar el usuario
	err := service.Delete(context.Background(), strconv.Itoa(created.ID))

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Intentar obtener el usuario eliminado
	_, err = service.GetByID(context.Background(), strconv.Itoa(created.ID))
	if err == nil {
		t.Error("Expected user to be deleted")
	}

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound after deletion, got %v", err)
	}
}
