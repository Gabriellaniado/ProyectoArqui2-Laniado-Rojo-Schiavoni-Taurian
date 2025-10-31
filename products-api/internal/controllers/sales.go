package controllers

import (
	"context"
	"errors"
	"net/http"
	"products-api/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

// SalesService define la lógica de negocio para Sales
// Capa intermedia entre Controllers (HTTP) y Repository (datos)
// Responsabilidades: validaciones, transformaciones, reglas de negocio
type SalesService interface {

	// Create valida y crea una nueva venta
	Create(ctx context.Context, sale domain.BodySales) (domain.Sales, error)

	// GetByID obtiene una venta por su ID
	GetByID(ctx context.Context, id string) (domain.Sales, error)

	// GetByCustomerID obtiene todas las ventas de un cliente
	GetByCustomerID(ctx context.Context, customerID string) ([]domain.Sales, error)

	// Update actualiza una venta existente
	Update(ctx context.Context, id string, sale domain.BodySales) (domain.Sales, error)

	// Delete elimina una venta por ID
	Delete(ctx context.Context, id string) error
}

// SalesController maneja las peticiones HTTP para Sales
// Responsabilidades:
// - Extraer datos del request (JSON, path params, query params)
// - Validar formato de entrada
// - Llamar al service correspondiente
// - Retornar respuesta HTTP adecuada
type SalesController struct {
	service SalesService // Inyección de dependencia
}

/*
const (
	salesDefaultPage  = 1
	salesDefaultCount = 10
)
*/

// NewSalesController crea una nueva instancia del controller
func NewSalesController(salesService SalesService) *SalesController {
	return &SalesController{
		service: salesService,
	}
}

// parseDate es una función helper para parsear fechas en diferentes formatos
func parseDate(dateStr string) (time.Time, error) {
	// Intentar formato RFC3339 (2024-01-01T00:00:00Z)
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}

	// Intentar formato de solo fecha (2024-01-01)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}

	// Intentar formato con hora (2024-01-01 15:04:05)
	if t, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("invalid date format")
}

// CreateSale maneja POST /sales - Crea una nueva venta
func (c *SalesController) CreateSale(ctx *gin.Context) {
	var sale domain.BodySales
	if err := ctx.ShouldBindJSON(&sale); err != nil {
		// Error en el formato del JSON
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	created, err := c.service.Create(ctx.Request.Context(), sale)
	if err != nil {
		// Error interno del servidor o validación
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create sale",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con la venta creada
	ctx.JSON(http.StatusCreated, gin.H{
		"sale": created,
	})
}

// GetSaleByID maneja GET /sales/:id - Obtiene venta por ID (MongoDB ObjectID)
func (c *SalesController) GetSaleByID(ctx *gin.Context) {
	sale, err := c.service.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "sale not found: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sale": sale,
	})
}

// GetSalesByCustomerID maneja GET /sales/customer/:customerID - Obtiene todas las ventas de un cliente
func (c *SalesController) GetSalesByCustomerID(ctx *gin.Context) {
	customerID := ctx.Param("customerID")
	if customerID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "customerID parameter is required",
		})
		return
	}

	sales, err := c.service.GetByCustomerID(ctx.Request.Context(), customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get sales by customer",
			"details": err.Error(),
		})
		return
	}

	// Calcular el total gastado por el cliente
	var totalSpent float64
	for _, sale := range sales {
		totalSpent += sale.TotalPrice
	}

	// Usar la estructura del domain
	response := domain.CustomerSalesResponse{
		CustomerID: customerID,
		Sales:      sales,
		Count:      len(sales),
		TotalSpent: totalSpent,
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateSale maneja PUT /sales/:id - Actualiza venta existente
func (c *SalesController) UpdateSale(ctx *gin.Context) {
	var updatedSale domain.BodySales

	if err := ctx.ShouldBindJSON(&updatedSale); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid JSON",
			"details": err.Error(),
		})
		return
	}

	sale, err := c.service.Update(ctx.Request.Context(), ctx.Param("id"), updatedSale)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update sale",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sale": sale,
	})
}

// DeleteSale maneja DELETE /sales/:id - Elimina venta por ID
func (c *SalesController) DeleteSale(ctx *gin.Context) {
	id := ctx.Param("id")

	// Validar que el ID no esté vacío
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID format",
		})
		return
	}

	err := c.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to delete sale",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "sale deleted successfully",
	})
}
