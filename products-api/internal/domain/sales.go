package domain

import (
	"time"
)

type Sales struct {
	ID         string    `json:"id"`
	ItemID     string    `json:"item_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	SaleDate   time.Time `json:"sale_date"`
	CustomerID string    `json:"customer_id"`
}

type BodySales struct {
	ItemID     string `json:"item_id"`
	Quantity   int    `json:"quantity"`
	CustomerID string `json:"customer_id"`
}

type SalesPaginatedResponse struct {
	Page    int     `json:"page"`
	Count   int     `json:"count"`
	Total   int     `json:"total"`
	Results []Sales `json:"results"`
}

type SalesSearchFilters struct {
	ID          string     `json:"id"`
	ItemID      string     `json:"item_id"`
	CustomerID  string     `json:"customer_id"`
	MinPrice    *float64   `json:"min_price"`
	MaxPrice    *float64   `json:"max_price"`
	MinQuantity *int       `json:"min_quantity"`
	MaxQuantity *int       `json:"max_quantity"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	SortBy      string     `json:"sort_by"`
	Page        int        `json:"page"`
	Count       int        `json:"count"`
}

type CustomerSalesResponse struct {
	CustomerID string  `json:"customer_id"`
	Sales      []Sales `json:"sales"`
	Count      int     `json:"count"`
	TotalSpent float64 `json:"total_spent"`
}
