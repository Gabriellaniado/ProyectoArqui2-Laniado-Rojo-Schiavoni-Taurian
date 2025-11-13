package repository

import (
	"context"
	"fmt"
	"products-api/internal/domain"
	"strconv"
	"time"

	"github.com/karlseguin/ccache"
)

// CartLocalCacheRepository implementa CartRepository usando CCache
type CartLocalCacheRepository struct {
	cache *ccache.Cache
	ttl   time.Duration
}

// NewCartLocalCacheRepository crea una nueva instancia del cache
func NewCartLocalCacheRepository(ttl time.Duration) *CartLocalCacheRepository {
	cache := ccache.New(ccache.Configure().MaxSize(1000).ItemsToPrune(100))
	return &CartLocalCacheRepository{
		cache: cache,
		ttl:   ttl,
	}
}

// GetByCustomerID obtiene un carrito del cache
func (r *CartLocalCacheRepository) GetByCustomerID(ctx context.Context, customerID int) (domain.Cart, error) {
	key := r.buildKey(customerID)
	item := r.cache.Get(key)
	if item == nil || item.Expired() {
		return domain.Cart{}, fmt.Errorf("cart not found in cache")
	}

	cart, ok := item.Value().(domain.Cart)
	if !ok {
		return domain.Cart{}, fmt.Errorf("invalid cart type in cache")
	}

	return cart, nil
}

// Create almacena un carrito en el cache
func (r *CartLocalCacheRepository) Create(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	key := r.buildKey(cart.CustomerID)
	r.cache.Set(key, cart, r.ttl)
	return cart, nil
}

// Update actualiza un carrito en el cache
func (r *CartLocalCacheRepository) Update(ctx context.Context, customerID int, cart domain.Cart) (domain.Cart, error) {
	key := r.buildKey(customerID)
	r.cache.Set(key, cart, r.ttl)
	return cart, nil
}

// Delete elimina un carrito del cache
func (r *CartLocalCacheRepository) Delete(ctx context.Context, customerID int) error {
	key := r.buildKey(customerID)
	r.cache.Delete(key)
	return nil
}

// Upsert inserta o actualiza un carrito en el cache
func (r *CartLocalCacheRepository) Upsert(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	key := r.buildKey(cart.CustomerID)
	r.cache.Set(key, cart, r.ttl)
	return cart, nil
}

// buildKey construye la clave del cache
func (r *CartLocalCacheRepository) buildKey(customerID int) string {
	return "cart:customer:" + strconv.Itoa(customerID)
}
