package domain

import (
	"time"
)

type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchFilters struct {
	Name     string   `json:"name"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	Category string   `json:"category"`
	SortBy   string   `json:"sort_by"`
	Page     int      `json:"page"`
	Count    int      `json:"count"`
}

type PaginatedResponse struct {
	Page    int    `json:"page"`
	Count   int    `json:"count"`
	Total   int    `json:"total"`
	Results []Item `json:"results"`
}
