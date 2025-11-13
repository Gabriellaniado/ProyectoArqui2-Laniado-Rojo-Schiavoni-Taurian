package main

import (
	"context"
	"log"
	"net/http"
	"products-api/internal/clients"
	"products-api/internal/config"
	"products-api/internal/controllers"
	"products-api/internal/middleware"
	"products-api/internal/repository"
	"products-api/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 📋 Cargar configuración desde las variables de entorno
	cfg := config.Load()

	// 🏗️ Inicializar capas de la aplicación (Dependency Injection)
	// Patrón: Repository -> Service -> Controller
	// Cada capa tiene una responsabilidad específica

	// Context
	ctx := context.Background()

	// Capa de datos: maneja operaciones DB
	itemsMongoRepo := repository.NewMongoItemsRepository(ctx, cfg.Mongo.URI, cfg.Mongo.DB, "items")

	// Capa de cache distribuida: maneja operaciones con Memcached
	itemsMemcachedRepo := repository.NewMemcachedItemsRepository(
		cfg.Memcached.Host,
		cfg.Memcached.Port,
		time.Duration(cfg.Memcached.TTLSeconds)*time.Second,
	)

	// Capa de cache local: maneja operaciones con CCache
	itemsLocalCacheRepo := repository.NewItemsLocalCacheRepository(30 * time.Second)

	// Inicializamos RabbitMQ para comunicar las novedades de escritura de items
	itemsQueue := clients.NewRabbitMQClient(
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.QueueName,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	// Capa de lógica de negocio: validaciones, transformaciones
	itemService := services.NewItemsService(itemsMongoRepo, itemsLocalCacheRepo, itemsMemcachedRepo, itemsQueue)

	// Capa de controladores: maneja HTTP requests/responses
	itemController := controllers.NewItemsController(&itemService)

	// ========================================
	// SALES - Configuración
	// ========================================

	// Repositorio MongoDB para Sales
	salesMongoRepo := repository.NewMongoSalesRepository(ctx, cfg.Mongo.URI, cfg.Mongo.DB, "sales")

	// Repositorio de cache local para Sales (Cache)
	salesLocalCacheRepo := repository.NewSalesLocalCacheRepository(30 * time.Second)

	// Capa de lógica de negocio para Sales (inyectamos itemService para calcular precios)
	salesService := services.NewSalesService(salesMongoRepo, salesLocalCacheRepo, &itemService)

	// Capa de controladores para Sales
	salesController := controllers.NewSalesController(&salesService)

	// Capa de lógica de negocio para Auth y controlador
	authService := services.NewAuthService("http://users-api:8082")
	authController := controllers.NewAuthController(authService)

	// 🌐 Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware: funciones que se ejecutan en cada request
	router.Use(middleware.CORSMiddleware)

	// 🏥 Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 📚 Rutas de Items API

	router.POST("/items", authController.VerifyAdminToken, itemController.CreateItem)

	// GET /items/:id - obtener item por ID
	router.GET("/items/:id", itemController.GetItemByID)

	// PUT /items/:id - actualizar item existente
	router.PUT("/items/:id", authController.VerifyAdminToken, itemController.UpdateItem)

	// DELETE /items/:id - eliminar item
	router.DELETE("/items/:id", authController.VerifyAdminToken, itemController.DeleteItem)

	// ========================================
	// SALES - Rutas
	// ========================================

	// GET /sales - listar ventas con filtros
	//router.GET("/sales", salesController.List)

	// POST /sales - crear nueva venta
	router.POST("/sales", authController.VerifyToken, salesController.CreateSale)

	// GET /sales/:id - obtener venta por ID (MongoDB ObjectID)
	router.GET("/sales/:id", authController.VerifyToken, salesController.GetSaleByID)

	// GET /sales/customer/:customerID - obtener todas las ventas de un cliente
	router.GET("/sales/customer/:customerID", authController.VerifyToken, salesController.GetSalesByCustomerID)

	// PUT /sales/:id - actualizar venta existente
	router.PUT("/sales/:id", authController.VerifyToken, salesController.UpdateSale)

	// DELETE /sales/:id - eliminar venta
	router.DELETE("/sales/:id", authController.VerifyToken, salesController.DeleteSale)

	// Configuración del server HTTP
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("📚 Items API: http://localhost:%s/items", cfg.Port)
	log.Printf("💰 Sales API: http://localhost:%s/sales", cfg.Port)

	// Iniciar servidor (bloquea hasta que se pare el servidor)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
