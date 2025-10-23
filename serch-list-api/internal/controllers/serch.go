package controllers

import (
	"context"
	"net/http"
	"serch-list-api/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SerchService interface {
	List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error)
}

type SerchController struct {
	service SerchService // Inyección de dependencia
}

func NewSerchController(SerchService SerchService) *SerchController {
	return &SerchController{
		service: SerchService,
	}
}

const (
	listDefaultPage  = 1
	listDefaultCount = 10
)

func (c *SerchController) List(ctx *gin.Context) {
	// Parsear filtros desde query params
	// Ejemplo GET /items?q=iphone&minPrice=100&maxPrice=500&page=2&count=20&sortBy=price%20desc
	filters := domain.SearchFilters{}

	filters.Name = ctx.Query("q")

	if minPriceStr := ctx.Query("minPrice"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := ctx.Query("maxPrice"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filters.Page = page
		}
	} else {
		filters.Page = listDefaultPage // default
	}

	if countStr := ctx.Query("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil {
			filters.Count = count
		}
	} else {
		filters.Count = listDefaultCount // default
	}

	filters.SortBy = ctx.DefaultQuery("sortBy", "createdAt desc")

	// 🔍 Llamar al service
	resp, err := c.service.List(ctx.Request.Context(), filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch items",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con paginación incluida
	ctx.JSON(http.StatusOK, resp)
}
