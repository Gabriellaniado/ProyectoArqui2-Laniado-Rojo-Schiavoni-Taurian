package main

import (
	"log"
	"net/http"
	"time"
	"users-api/internal/config"
	"users-api/internal/controllers"
	"users-api/internal/db"
	"users-api/internal/middleware"
	"users-api/internal/repository"
	"users-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 📋 Cargar configuración desde las variables de entorno
	cfg := config.Load()

	// 🔌 Conectar a MySQL con GORM
	// Establecemos conexión y obtenemos referencia a la base de datos
	mysqlDB, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}

	// 🏗️ Inicializar capas de la aplicación (Dependency Injection)
	// Patrón: Repository -> Service -> Controller
	// Cada capa tiene una responsabilidad específica

	// Capa de datos: maneja operaciones MySQL con GORM
	userRepo := repository.NewMySQLUsersRepository(mysqlDB)

	// Capa de lógica de negocio: validaciones, transformaciones
	userService := services.NewUsersService(userRepo)

	// Capa de controladores: maneja HTTP requests/responses
	userController := controllers.NewUsersController(userService)

	// 🌐 Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware: funciones que se ejecutan en cada request
	router.Use(middleware.CORSMiddleware)

	// 🏥 Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 📚 Rutas de Users API
	// GET /users - listar todos los usuarios
	router.GET("/users", userController.GetUsers)

	// POST /users - crear nuevo usuario
	router.POST("/users", userController.CreateUser)

	// GET /users/:id - obtener usuario por ID
	router.GET("/users/:id", userController.GetUserByID)

	// GET /users/email/:email - obtener usuario por email
	router.GET("/users/email/:email", userController.GetUserByEmail)

	// PUT /users/:id - actualizar usuario existente
	router.PUT("/users/:id", userController.UpdateUser)

	// DELETE /users/:id - eliminar usuario
	router.DELETE("/users/:id", userController.DeleteUser)

	// POST /auth/login - login de usuario
	router.POST("/auth/login", userController.Login)

	router.POST("/auth/verify-token", userController.VerifyToken)

	router.POST("/auth/verify-admin-token", userController.VerifyAdminToken)

	// Configuración del server HTTP
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("� Users API: http://localhost:%s/users", cfg.Port)

	// Iniciar servidor (bloquea hasta que se pare el servidor)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
