package repository

import (
	"context"
	"fmt"
	"products-api/internal/domain"
	"time"

	"github.com/karlseguin/ccache"
)

type SalesLocalCacheRepository struct {
	client *ccache.Cache
	ttl    time.Duration
}

func NewSalesLocalCacheRepository(ttl time.Duration) *SalesLocalCacheRepository {
	return &SalesLocalCacheRepository{
		client: ccache.New(ccache.Configure()),
		ttl:    ttl,
	}
}

func (r SalesLocalCacheRepository) Create(ctx context.Context, sale domain.Sales) (domain.Sales, error) {
	r.client.Set(sale.ID, sale, r.ttl)
	return sale, nil
}

func (r SalesLocalCacheRepository) GetByID(ctx context.Context, id string) (domain.Sales, error) {
	it := r.client.Get(id)
	if it == nil {
		return domain.Sales{}, fmt.Errorf("sale not found in cache")
	}
	sale, ok := it.Value().(domain.Sales)
	if !ok {
		return domain.Sales{}, fmt.Errorf("error asserting sale type from cache")
	}
	return sale, nil
}

func (r SalesLocalCacheRepository) GetByCustomerID(ctx context.Context, customerID int) ([]domain.Sales, error) {
	// Cache local solo puede buscar por clave exacta (ID)
	return []domain.Sales{}, fmt.Errorf("getByCustomerID is not supported in local cache")
}

func (r SalesLocalCacheRepository) Update(ctx context.Context, id string, sale domain.Sales) (domain.Sales, error) {
	r.client.Set(id, sale, r.ttl)
	return sale, nil
}

func (r SalesLocalCacheRepository) Delete(ctx context.Context, id string) error {
	it := r.client.Delete(id)
	if !it {
		return fmt.Errorf("sale not found in cache")
	}
	return nil
}
