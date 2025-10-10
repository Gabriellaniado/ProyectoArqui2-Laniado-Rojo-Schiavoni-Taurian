package controllers

import (
	"clase02-mongo/internal/domain"
	"clase02-mongo/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ItemsController maneja las peticiones HTTP para Items
// Responsabilidades:
// - Extraer datos del request (JSON, path params, query params)
// - Validar formato de entrada
// - Llamar al service correspondiente
// - Retornar respuesta HTTP adecuada
type ItemsController struct {
	service services.ItemsService // Inyección de dependencia
}

// NewItemsController crea una nueva instancia del controller
func NewItemsController(itemsService services.ItemsService) *ItemsController {
	return &ItemsController{
		service: itemsService,
	}
}

// GetItems maneja GET /items - Lista todos los items
// ✅ IMPLEMENTADO - Ejemplo para que los estudiantes entiendan el patrón
func (c *ItemsController) GetItems(ctx *gin.Context) {
	// 🔍 Llamar al service para obtener los datos
	items, err := c.service.List(ctx.Request.Context())
	if err != nil {
		// ❌ Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch items",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con los datos
	ctx.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// CreateItem maneja POST /items - Crea un nuevo item
// Consigna 1: Recibir JSON, validar y crear item
func (c *ItemsController) CreateItem(ctx *gin.Context) {
	// Consigna 1: Recibir JSON, validar y crear item

	var newItem struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	// Bindear JSON al struct
	if err := ctx.ShouldBindJSON(&newItem); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Llamar al service para crear el nuevo item
	createdItem, err := c.service.Create(ctx.Request.Context(), domain.Item{
		Name:  newItem.Name,
		Price: newItem.Price,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create item",
			"details": err.Error(),
		})
		return
	}

	// Respuesta exitosa con el item creado
	ctx.JSON(http.StatusCreated, gin.H{
		"item": createdItem,
	})
}

// GetItemByID maneja GET /items/:id - Obtiene item por ID
// Consigna 2: Extraer ID del path param, validar y buscar
func (c *ItemsController) GetItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id parameter is required",
		})
		return
	}
	item, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch item",
			"details": err.Error(),
		})
		return
	}
	if item.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Item not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"item": item,
	})
}

// UpdateItem maneja PUT /items/:id - Actualiza item existente
// Consigna 3: Extraer ID y datos, validar y actualizar
func (c *ItemsController) UpdateItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id parameter is required",
		})
		return
	}
	var updateData struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}
	updatedItem, err := c.service.Update(ctx.Request.Context(), id, domain.Item{
		Name:  updateData.Name,
		Price: updateData.Price,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update item",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"item": updatedItem,
	})
}

// DeleteItem maneja DELETE /items/:id - Elimina item por ID
// Consigna 4: Extraer ID, validar y eliminar
func (c *ItemsController) DeleteItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id parameter is required",
		})
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete item",
			"details": err.Error(),
		})
		return
	}
	ctx.Status(http.StatusNoContent) // 204 No Content

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
