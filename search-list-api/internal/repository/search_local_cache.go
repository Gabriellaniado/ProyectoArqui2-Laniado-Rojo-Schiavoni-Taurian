package repository

import (
	"context"
	"errors"
	"search-list-api/internal/domain"
	"time"

	"github.com/karlseguin/ccache"
)

type SearchLocalCacheRepository struct {
	client *ccache.Cache
	ttl    time.Duration
}

func NewSearchLocalCacheRepository(ttl time.Duration) *SearchLocalCacheRepository {
	return &SearchLocalCacheRepository{
		client: ccache.New(ccache.Configure()),
		ttl:    ttl,
	}
}

func (r SearchLocalCacheRepository) ListHash(ctx context.Context, hash string) (domain.PaginatedResponse, error) {
	item := r.client.Get(hash)
	if item == nil || item.Expired() {
		return domain.PaginatedResponse{}, errors.New("cache miss")
	}

	result, ok := item.Value().(domain.PaginatedResponse)
	if !ok {
		return domain.PaginatedResponse{}, errors.New("invalid cache data type")
	}

	return result, nil
}

func (r SearchLocalCacheRepository) SaveWithHash(ctx context.Context, hash string, response domain.PaginatedResponse) error {
	r.client.Set(hash, response, r.ttl)
	return nil
}
