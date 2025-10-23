package controllers

import (
	"context"
	"net/http"
	"products-api/internal/domain"

	"github.com/gin-gonic/gin"
)

// ItemsService define la lógica de negocio para Items
// Capa intermedia entre Controllers (HTTP) y Repository (datos)
// Responsabilidades: validaciones, transformaciones, reglas de negocio
type ItemsService interface {

	// Create valida y crea un nuevo item
	Create(ctx context.Context, item domain.Item) (domain.Item, error)

	// GetByID obtiene un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)

	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)

	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
}

// ItemsController maneja las peticiones HTTP para Items
// Responsabilidades:
// - Extraer datos del request (JSON, path params, query params)
// - Validar formato de entrada
// - Llamar al service correspondiente
// - Retornar respuesta HTTP adecuada
type ItemsController struct {
	service ItemsService // Inyección de dependencia
}

// NewItemsController crea una nueva instancia del controller
func NewItemsController(itemsService ItemsService) *ItemsController {
	return &ItemsController{
		service: itemsService,
	}
}

// CreateItem maneja POST /items - Crea un nuevo item
// Consigna 1: Recibir JSON, validar y crear item
func (c *ItemsController) CreateItem(ctx *gin.Context) {
	var item domain.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		// ❌ Error en el formato del JSON
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	created, err := c.service.Create(ctx.Request.Context(), item)
	if err != nil {
		// ❌ Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create item",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con el item creado
	ctx.JSON(http.StatusCreated, gin.H{
		"item": created,
	})
}

// GetItemByID maneja GET /items/:id - Obtiene item por ID
// Consigna 2: Extraer ID del path param, validar y buscar
func (c *ItemsController) GetItemByID(ctx *gin.Context) {
	item, err := c.service.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "item not found: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"item": item})
}

// UpdateItem maneja PUT /items/:id - Actualiza item existente
// Consigna 3: Extraer ID y datos, validar y actualizar
func (c *ItemsController) UpdateItem(ctx *gin.Context) {
	var updatedItem domain.Item

	if err := ctx.ShouldBindJSON(&updatedItem); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	item, err := c.service.Update(ctx.Request.Context(), ctx.Param("id"), updatedItem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update item: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"item": item})
}

// DeleteItem maneja DELETE /items/:id - Elimina item por ID
// Consigna 4: Extraer ID, validar y eliminar
func (c *ItemsController) DeleteItem(ctx *gin.Context) {
	id := ctx.Param("id")
	//valida id
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	err := c.service.Delete(ctx.Request.Context(), id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}

// 📚 Notas sobre HTTP Status Codes
//
// 200 OK - Operación exitosa con contenido
// 201 Created - Recurso creado exitosamente
// 204 No Content - Operación exitosa sin contenido (típico para DELETE)
// 400 Bad Request - Error en los datos enviados por el cliente
// 404 Not Found - Recurso no encontrado
// 500 Internal Server Error - Error interno del servidor
// 501 Not Implemented - Funcionalidad no implementada (para TODOs)
//
// 💡 Tip: En una API real, sería buena práctica crear una función
// helper para manejar respuestas de error de manera consistente
