package repository

import (
	"context"
	"fmt"
	"search-list-api/internal/clients"
	"search-list-api/internal/domain"
	"strings"
)

type SolrClient interface {
}

// SolrItemsRepository implementa ItemsRepository usando Solr
type SolrItemsRepository struct {
	client *clients.SolrClient
}

// NewSolrItemsRepository crea una nueva instancia del repository
func NewSolrItemsRepository(host, port, core string) *SolrItemsRepository {
	client := clients.NewSolrClient(host, port, core)
	return &SolrItemsRepository{
		client: client,
	}
}

// List retorna items desde Solr en base a los filtros
func (r *SolrItemsRepository) List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error) {
	query := buildQuery(filters)

	return r.client.Search(ctx, query, filters.Page, filters.Count)
}

// Create indexa un nuevo item en Solr
func (r *SolrItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	if err := r.client.Index(ctx, item); err != nil {
		return domain.Item{}, fmt.Errorf("error indexing item in solr: %w", err)
	}
	return item, nil
}

// Update actualiza un item existente en Solr (re-indexa)
func (r *SolrItemsRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	// En Solr, actualizar es equivalente a re-indexar con el mismo ID
	item.ID = id
	if err := r.client.Index(ctx, item); err != nil {
		return domain.Item{}, fmt.Errorf("error updating item in solr: %w", err)
	}
	return item, nil
}

// Delete elimina un item por ID de Solr
func (r *SolrItemsRepository) Delete(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting item from solr: %w", err)
	}
	return nil
}

// buildQuery construye la query de Solr a partir de los filtros
func buildQuery(filters domain.SearchFilters) string {
	var parts []string

	if filters.Name != "" {
		// 1. Separamos el string de búsqueda por espacios
		// Ej: "mates argentinos" -> ["mates", "argentinos"]
		terms := strings.Fields(filters.Name)

		var nameParts []string
		for _, term := range terms {
			nameParts = append(nameParts, fmt.Sprintf("name:%s~1", term))
		}

		// 3. Unimos los términos con AND y los agrupamos
		// Resultado: (name:mates~1 AND name:argentinos~1)
		if len(nameParts) > 0 {
			parts = append(parts, "("+strings.Join(nameParts, " AND ")+")")
		}
	}
	// Filtro por categoría
	if filters.Category != "" {
		parts = append(parts, fmt.Sprintf("category:%s", filters.Category))
	}

	// Filtro por rango de precios
	if filters.MinPrice != nil || filters.MaxPrice != nil {
		minStr := "*"
		maxStr := "*"

		if filters.MinPrice != nil {
			// Usamos %g para un formato de float más limpio (ej. 100)
			// en lugar de %f (ej. 100.000000)
			minStr = fmt.Sprintf("%g", *filters.MinPrice)
		}
		if filters.MaxPrice != nil {
			maxStr = fmt.Sprintf("%g", *filters.MaxPrice)
		}

		parts = append(parts, fmt.Sprintf("price:[%s TO %s]", minStr, maxStr))
	}

	if len(parts) == 0 {
		return "*:*" // Query que retorna todo si no hay filtros
	}

	return strings.Join(parts, " AND ")
}
