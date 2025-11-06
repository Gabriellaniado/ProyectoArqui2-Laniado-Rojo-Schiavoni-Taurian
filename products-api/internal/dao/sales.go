package dao

import (
	"products-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sales struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ItemID     string             `bson:"item_id"`
	Quantity   int                `bson:"quantity"`
	TotalPrice float64            `bson:"total_price"`
	SaleDate   time.Time          `bson:"sale_date"`
	CustomerID int                `bson:"customer_id"`
}

type SalesList []Sales

func (sl SalesList) ToDomainList() []domain.Sales {
	var result []domain.Sales
	for _, s := range sl {
		result = append(result, s.ToDomain())
	}
	return result
}

func (s Sales) ToDomain() domain.Sales {
	return domain.Sales{
		ID:         s.ID.Hex(),
		ItemID:     s.ItemID,
		Quantity:   s.Quantity,
		TotalPrice: s.TotalPrice,
		SaleDate:   s.SaleDate,
		CustomerID: s.CustomerID,
	}
}

func (s Sales) TimeInfo() (time.Time, time.Time) {
	return s.SaleDate, time.Time{}
}

func FromDomainSales(domainSales domain.Sales) Sales {
	var objectID primitive.ObjectID
	if domainSales.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(domainSales.ID)
	}
	return Sales{
		ID:         objectID,
		ItemID:     domainSales.ItemID,
		Quantity:   domainSales.Quantity,
		TotalPrice: domainSales.TotalPrice,
		SaleDate:   domainSales.SaleDate,
		CustomerID: domainSales.CustomerID,
	}
}
