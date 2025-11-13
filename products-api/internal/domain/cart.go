package domain

import (
	"time"
)

// CartItem representa un ítem individual dentro del carrito
type CartItem struct {
	ItemID   string `json:"item_id" bson:"item_id"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

// Cart representa el carrito de compras de un usuario
type Cart struct {
	ID         string     `json:"id" bson:"_id,omitempty"`
	CustomerID int        `json:"customer_id" bson:"customer_id"`
	Items      []CartItem `json:"items" bson:"items"`
	Total      float64    `json:"total" bson:"total"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
}

// AddItemRequest representa la request para agregar un ítem al carrito
type AddItemRequest struct {
	ItemID   string `json:"item_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

// UpdateItemRequest representa la request para actualizar un ítem del carrito
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// CartResponse representa la respuesta del carrito con información enriquecida
type CartResponse struct {
	ID         string                `json:"id"`
	CustomerID int                   `json:"customer_id"`
	Items      []CartItemWithDetails `json:"items"`
	Total      float64               `json:"total"`
	ItemCount  int                   `json:"item_count"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

// CartItemWithDetails incluye la información completa del producto
type CartItemWithDetails struct {
	ItemID      string  `json:"item_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
	Stock       int     `json:"stock"` // Stock disponible del producto
}

// CheckoutRequest representa la request para finalizar una compra
type CheckoutRequest struct {
	// Podemos agregar más campos en el futuro (dirección, método de pago, etc.)
}
