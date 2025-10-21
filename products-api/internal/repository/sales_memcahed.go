package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"products-api/internal/domain"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedSalesRepository struct {
	ttl    time.Duration
	client *memcache.Client
}

func NewMemcachedSalesRepository(host string, port string, ttl time.Duration) MemcachedSalesRepository {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))

	return MemcachedSalesRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r MemcachedSalesRepository) Create(ctx context.Context, sale domain.Sales) (domain.Sales, error) {
	bytes, err := json.Marshal(sale)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error marshalling sale to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        sale.ID,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return domain.Sales{}, fmt.Errorf("error setting sale in memcached: %w", err)
	}
	return sale, nil
}

func (r MemcachedSalesRepository) GetByID(ctx context.Context, id string) (domain.Sales, error) {
	bytes, err := r.client.Get(id)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error getting sale from memcached: %w", err)
	}
	var sale domain.Sales
	if err := json.Unmarshal(bytes.Value, &sale); err != nil {
		return domain.Sales{}, fmt.Errorf("error unmarshalling sale from JSON: %w", err)
	}
	return sale, nil
}

func (r MemcachedSalesRepository) GetByCustomerID(ctx context.Context, customerID string) ([]domain.Sales, error) {
	// Memcached solo soporta búsqueda por clave exacta (ID)
	// No podemos buscar por CustomerID sin indexación adicional
	return []domain.Sales{}, fmt.Errorf("getByCustomerID is not supported in memcached")
}

func (r MemcachedSalesRepository) Update(ctx context.Context, id string, sale domain.Sales) (domain.Sales, error) {
	bytes, err := json.Marshal(sale)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error marshalling sale to JSON: %w", err)
	}

	if err := r.client.Replace(&memcache.Item{
		Key:        id,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return domain.Sales{}, fmt.Errorf("error replacing sale in memcached: %w", err)
	}
	return sale, nil
}

func (r MemcachedSalesRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Delete(id)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			// Sale doesn't exist in cache, which is fine - return nil
			return nil
		}
		return fmt.Errorf("error deleting sale from memcached: %w", err)
	}
	return nil
}
