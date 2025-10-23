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
	// itemsLocalCacheRepo := repository.NewItemsLocalCacheRepository(30 * time.Second)

	// Capa de búsqueda: maneja operaciones de búsqueda con Solr
	/*itemsSolrRepo := repository.NewSolrItemsRepository(
		cfg.Solr.Host,
		cfg.Solr.Port,
		cfg.Solr.Core,
	)*/

	// Inicializamos RabbitMQ para comunicar las novedades de escritura de items
	itemsQueue := clients.NewRabbitMQClient(
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.QueueName,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	// Capa de lógica de negocio: validaciones, transformaciones
	itemService := services.NewItemsService(itemsMongoRepo, itemsMemcachedRepo, itemsQueue, itemsQueue)
	//go itemService.InitConsumer(ctx)

	// Capa de controladores: maneja HTTP requests/responses
	itemController := controllers.NewItemsController(&itemService)

	// Cache (ejercicio: ajustar TTL y agregar "índice" de claves)
	// cache := cache.NewMemcached(memAddr)

	// 🌐 Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware: funciones que se ejecutan en cada request
	router.Use(middleware.CORSMiddleware)

	// 🏥 Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 📚 Rutas de Items API
	// GET /items - listar los items con filtros(✅ implementado)
	//router.GET("/items", itemController.List)

	router.POST("/items", itemController.CreateItem)

	// GET /items/:id - obtener item por ID
	router.GET("/items/:id", itemController.GetItemByID)

	// PUT /items/:id - actualizar item existente
	router.PUT("/items/:id", itemController.UpdateItem)

	// DELETE /items/:id - eliminar item
	router.DELETE("/items/:id", itemController.DeleteItem)

	// Configuración del server HTTP
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("📚 Items API: http://localhost:%s/items", cfg.Port)

	// Iniciar servidor (bloquea hasta que se pare el servidor)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
