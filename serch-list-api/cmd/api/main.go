package main

import (
	"context"
	"log"
	"net/http"
	"serch-list-api/internal/clients"
	"serch-list-api/internal/config"
	"serch-list-api/internal/controllers"
	"serch-list-api/internal/middleware"
	"serch-list-api/internal/repository"
	"serch-list-api/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 📋 Cargar configuración desde las variables de entorno
	cfg := config.Load()

	// Context
	ctx := context.Background()

	// Capa de búsqueda: maneja operaciones de búsqueda con Solr
	itemsSolrRepo := repository.NewSolrItemsRepository(
		cfg.Solr.Host,
		cfg.Solr.Port,
		cfg.Solr.Core,
	)

	// Inicializamos RabbitMQ para comunicar las novedades de escritura de items
	itemsQueue := clients.NewRabbitMQClient(
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.QueueName,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	// Capa de lógica de negocio: validaciones, transformaciones
	SerchService := services.NewSerchService(itemsSolrRepo, itemsQueue)
	go SerchService.InitConsumer(ctx)

	// Capa de controladores: maneja HTTP requests/responses
	SerchController := controllers.NewSerchController(SerchService)

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
	router.GET("/items", SerchController.List)

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
