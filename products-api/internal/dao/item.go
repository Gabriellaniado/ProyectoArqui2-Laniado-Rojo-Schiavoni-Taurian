package dao

import (
	"products-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Category    string             `bson:"category"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Stock       int                `bson:"stock"`
	ImageURL    string             `bson:"image_url"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func (i Item) ToDomain() domain.Item {
	return domain.Item{
		ID:          i.ID.Hex(),
		Name:        i.Name,
		Category:    i.Category,
		Description: i.Description,
		Price:       i.Price,
		Stock:       i.Stock,
		ImageURL:    i.ImageURL,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func (i Item) TimeInfo() (time.Time, time.Time) {
	return i.CreatedAt, i.UpdatedAt
}

func FromDomain(domainItem domain.Item) Item {
	var objectID primitive.ObjectID
	if domainItem.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(domainItem.ID)
	}
	return Item{
		ID:          objectID,
		Name:        domainItem.Name,
		Category:    domainItem.Category,
		Description: domainItem.Description,
		Price:       domainItem.Price,
		Stock:       domainItem.Stock,
		ImageURL:    domainItem.ImageURL,
		CreatedAt:   domainItem.CreatedAt,
		UpdatedAt:   domainItem.UpdatedAt,
	}
}
