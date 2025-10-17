package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"users-api/internal/domain"

	"github.com/gin-gonic/gin"
)

// UsersController maneja las peticiones HTTP relacionadas con usuarios
// Actúa como intermediario entre las rutas HTTP y la lógica de negocio (services)

type UsersService interface {
	List(ctx context.Context) ([]domain.UserResponse, error)
	Create(ctx context.Context, user domain.User) (domain.UserResponse, error)
	GetByID(ctx context.Context, id string) (domain.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, id string, user domain.User) (domain.UserResponse, error)
	Delete(ctx context.Context, id string) error
	Login(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error)
}

type UsersController struct {
	service UsersService // Referencia al service para lógica de negocio
}

// NewUsersController crea una nueva instancia del controller
// Constructor que inyecta la dependencia del service
func NewUsersController(service UsersService) *UsersController {
	return &UsersController{
		service: service,
	}
}

// GetUsers obtiene todos los usuarios
// GET /users
func (c *UsersController) GetUsers(ctx *gin.Context) {
	// 1. Llamamos al service para obtener todos los usuarios
	users, err := c.service.List(ctx)
	if err != nil {
		// Si hay error, retornamos 500 (Internal Server Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*for i := range users {
		users[i].Password = "" // No devolver el hash de la contraseña HAY QUE BUSCAR OTRA FORMA, QUIZAS OTRO DTO.
	}*/
	// 2. Si todo sale bien, retornamos 200 con la lista de usuarios
	ctx.JSON(http.StatusOK, users)
}

// CreateUser crea un nuevo usuario
// POST /users
func (c *UsersController) CreateUser(ctx *gin.Context) {
	var user domain.User

	// 1. Parseamos el JSON del body de la request al struct User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		// Si el JSON es inválido, retornamos 400 (Bad Request)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Llamamos al service para crear el usuario
	createdUser, err := c.service.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "email") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		// Si hay error en la creación, retornamos 500
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//createdUser.Password = "" // No devolver el hash de la contraseña

	// 3. Si se crea exitosamente, retornamos 201 (Created) con el usuario
	ctx.JSON(http.StatusCreated, createdUser)
}

// GetUserByID obtiene un usuario por su ID
// GET /users/:id
func (c *UsersController) GetUserByID(ctx *gin.Context) {
	// 1. Extraemos el parámetro "id" de la URL
	idStr := ctx.Param("id")

	// 2. Convertimos el string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Si no es un número válido, retornamos 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 3. Buscamos el usuario por ID
	user, err := c.service.GetByID(ctx, strconv.Itoa(id))
	if err != nil {
		// Si no se encuentra o hay error, retornamos 404 (Not Found)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//user.Password = "" // No devolver el hash de la contraseña
	// 4. Si se encuentra, retornamos 200 con el usuario
	ctx.JSON(http.StatusOK, user)
}

// GetUserByEmail obtiene un usuario por su email
// GET /users/email/:email
func (c *UsersController) GetUserByEmail(ctx *gin.Context) {
	// 1. Extraemos el parámetro "email" de la URL
	email := ctx.Param("email")

	// 2. Buscamos el usuario por email
	user, err := c.service.GetByEmail(ctx, email)
	if err != nil {
		// Si no se encuentra o hay error, retornamos 404
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user.Password = "" // No devolver el hash de la contraseña. Es necesario que llegue la contraseña para el login

	// 3. Si se encuentra, retornamos 200 con el usuario
	ctx.JSON(http.StatusOK, user)
}

// UpdateUser actualiza un usuario existente
// PUT /users/:id
func (c *UsersController) UpdateUser(ctx *gin.Context) {
	// 1. Extraemos y validamos el ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 2. Parseamos el JSON del body
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Actualizamos el usuario
	updatedUser, err := c.service.Update(ctx, strconv.Itoa(id), user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//updatedUser.Password = "" // No devolver el hash de la contraseña
	// 4. Retornamos el usuario actualizado
	ctx.JSON(http.StatusOK, updatedUser)
}

// DeleteUser elimina un usuario por ID
// DELETE /users/:id
func (c *UsersController) DeleteUser(ctx *gin.Context) {
	// 1. Extraemos y validamos el ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 2. Eliminamos el usuario
	err = c.service.Delete(ctx, strconv.Itoa(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Retornamos mensaje de éxito
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (c *UsersController) Login(ctx *gin.Context) {
	var request domain.LoginRequest

	// 1. Recibir usuario y contraseña desde el body de la request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 2. Llamar al servicio de login (pasando el LoginRequest completo)
	response, err := c.service.Login(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No se pudo iniciar sesion"})
		return
	}

	// 3. Si el login es exitoso, devolver la respuesta
	ctx.JSON(http.StatusOK, response)
}
