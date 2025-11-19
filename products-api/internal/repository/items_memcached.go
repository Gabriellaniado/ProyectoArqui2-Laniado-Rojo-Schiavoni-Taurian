package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"products-api/internal/domain"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedItemsRepository struct {
	ttl    time.Duration
	client *memcache.Client
}

func NewMemcachedItemsRepository(host string, port string, ttl time.Duration) *MemcachedItemsRepository {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))

	return &MemcachedItemsRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *MemcachedItemsRepository) List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error) {
	return domain.PaginatedResponse{}, fmt.Errorf("list is not supported in memcached")
}

func (r *MemcachedItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	bytes, err := json.Marshal(item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error marshalling item to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        item.ID,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return domain.Item{}, fmt.Errorf("error setting item in memcached: %w", err)
	}
	return item, nil
}

func (r *MemcachedItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	bytes, err := r.client.Get(id)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error getting item from memcached: %w", err)
	}
	var item domain.Item
	if err := json.Unmarshal(bytes.Value, &item); err != nil {
		return domain.Item{}, fmt.Errorf("error unmarshalling item from JSON: %w", err)
	}
	if item.ID == "" {
		item.ID = id
	}
	return item, nil
}

func (r *MemcachedItemsRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	item.ID = id

	bytes, err := json.Marshal(item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error marshalling item to JSON: %w", err)
	}

	if _, err := r.client.Get(id); err == nil {

		if err := r.client.Replace(&memcache.Item{
			Key:        id,
			Value:      bytes,
			Expiration: int32(r.ttl.Seconds()),
		}); err != nil {
			return domain.Item{}, fmt.Errorf("error replacing item in memcached: %w", err)
		}
	}
	return item, nil
}

func (r *MemcachedItemsRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Delete(id)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			// Item doesn't exist in cache, which is fine - return nil
			return nil
		}
		return fmt.Errorf("error deleting item from memcached: %w", err)
	}
	return nil
}

func (r *MemcachedItemsRepository) DecrementStockAtomic(ctx context.Context, itemID string, quantity int) (bool, error) {
	// Memcached es solo cache, no maneja stock
	return false, nil
}

func (r *MemcachedItemsRepository) IncrementStock(ctx context.Context, itemID string, quantity int) error {
	// Memcached es solo cache, no maneja stock
	return nil
}
