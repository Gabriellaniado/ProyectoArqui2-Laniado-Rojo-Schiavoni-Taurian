package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItemDAO representa un ítem del carrito en la base de datos
type CartItemDAO struct {
	ItemID   string `bson:"item_id"`
	Quantity int    `bson:"quantity"`
}

// CartDAO representa el carrito en la base de datos MongoDB
type CartDAO struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID int                `bson:"customer_id"`
	Items      []CartItemDAO      `bson:"items"`
	Total      float64            `bson:"total"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}
