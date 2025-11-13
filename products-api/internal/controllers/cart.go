package controllers

import (
	"context"
	"log"
	"net/http"
	"products-api/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CartService define las operaciones de negocio para Cart
type CartService interface {
	GetCart(ctx context.Context, customerID int) (domain.CartResponse, error)
	AddItem(ctx context.Context, customerID int, req domain.AddItemRequest) (domain.CartResponse, error)
	UpdateItemCart(ctx context.Context, customerID int, itemID string, req domain.UpdateItemRequest) (domain.CartResponse, error)
	RemoveItem(ctx context.Context, customerID int, itemID string) (domain.CartResponse, error)
	ClearCart(ctx context.Context, customerID int) error
	Checkout(ctx context.Context, customerID int) ([]domain.Sales, error)
}

// CartController maneja las peticiones HTTP relacionadas con carritos
type CartController struct {
	service CartService
}

// NewCartController crea una nueva instancia del controller
func NewCartController(service CartService) *CartController {
	return &CartController{
		service: service,
	}
}

// GetCart obtiene el carrito del usuario autenticado
// GET /cart/:customerID
func (c *CartController) GetCart(ctx *gin.Context) {
	customerIDStr := ctx.Param("customerID")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer_id format",
		})
		return
	}

	cart, err := c.service.GetCart(ctx, customerID)
	if err != nil {
		log.Printf("❌ Error getting cart: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error retrieving cart",
		})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// AddItem agrega un item al carrito
// POST /cart/:customerID/items
func (ctrl *CartController) AddItem(c *gin.Context) {
	var req domain.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID, _ := strconv.Atoi(c.Param("customerID"))

	// Pasar c.Request.Context() en lugar de c
	cart, err := ctrl.service.AddItem(c.Request.Context(), customerID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateItem actualiza la cantidad de un item en el carrito
// PUT /cart/:customerID/items/:itemID
func (c *CartController) UpdateItemCart(ctx *gin.Context) {
	customerIDStr := ctx.Param("customerID")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer_id format",
		})
		return
	}

	itemID := ctx.Param("itemID")
	if itemID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "item_id is required",
		})
		return
	}

	var req domain.UpdateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	cart, err := c.service.UpdateItemCart(ctx, customerID, itemID, req)
	if err != nil {
		log.Printf("Error updating item in cart: %v", err)

		if err.Error() == "item not found in cart" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Item not found in cart",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// RemoveItem elimina un item del carrito
// DELETE /cart/:customerID/items/:itemID
func (c *CartController) RemoveItem(ctx *gin.Context) {
	customerIDStr := ctx.Param("customerID")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer_id format",
		})
		return
	}

	itemID := ctx.Param("itemID")
	if itemID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "item_id is required",
		})
		return
	}

	cart, err := c.service.RemoveItem(ctx, customerID, itemID)
	if err != nil {
		log.Printf(" Error removing item from cart: %v", err)

		if err.Error() == "item not found in cart" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Item not found in cart",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// ClearCart vacía el carrito completamente
// DELETE /cart/:customerID
func (c *CartController) ClearCart(ctx *gin.Context) {
	customerIDStr := ctx.Param("customerID")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer_id format",
		})
		return
	}

	err = c.service.ClearCart(ctx, customerID)
	if err != nil {
		log.Printf("❌ Error clearing cart: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error clearing cart",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cart cleared successfully",
	})
}

// Checkout procesa la compra del carrito
// POST /cart/:customerID/checkout
func (c *CartController) Checkout(ctx *gin.Context) {
	customerIDStr := ctx.Param("customerID")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer_id format",
		})
		return
	}

	sales, err := c.service.Checkout(ctx, customerID)
	if err != nil {
		log.Printf("❌ Error processing checkout: %v", err)

		if err.Error() == "cart is empty" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot checkout an empty cart",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Checkout completed successfully",
		"sales":   sales,
	})
}
