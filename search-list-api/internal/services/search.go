package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"search-list-api/internal/domain"
)

type SearchRepository interface {
	List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error)
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)
	Delete(ctx context.Context, id string) error
}

type SearchLocalCacheRepository interface {
	ListHash(ctx context.Context, hash string) (domain.PaginatedResponse, error)
	SaveWithHash(ctx context.Context, hash string, response domain.PaginatedResponse) error
}

type ItemsConsumer interface { //consume mensajes de rabbit
	Consume(ctx context.Context, handler func(ctx context.Context, message ItemEvent) error) error
}

type SearchServiceImpl struct {
	repo       SearchRepository
	consumer   ItemsConsumer
	localCache SearchLocalCacheRepository
}

func NewSearchService(repo SearchRepository, consumer ItemsConsumer, localCache SearchLocalCacheRepository) *SearchServiceImpl {
	return &SearchServiceImpl{
		repo:       repo,
		consumer:   consumer,
		localCache: localCache,
	}
}

func (s *SearchServiceImpl) List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error) {

	// Construir string con todos los filtros para generar hash único
	filterString := fmt.Sprintf("name:%s|minPrice:%v|maxPrice:%v|category:%s|sortBy:%s|page:%d|count:%d",
		filters.Name,
		filters.MinPrice,
		filters.MaxPrice,
		filters.Category,
		filters.SortBy,
		filters.Page,
		filters.Count,
	)

	hashBytes := md5.Sum([]byte(filterString))
	filterHash := hex.EncodeToString(hashBytes[:])

	// 1. Intentar obtener de caché
	cachedResult, err := s.localCache.ListHash(ctx, filterHash)
	if err == nil {
		return cachedResult, nil
	}

	// 2. Si no está en caché, obtener del repositorio (Solr)
	result, err := s.repo.List(ctx, filters)
	if err != nil {
		return domain.PaginatedResponse{}, fmt.Errorf("error listing from repository: %w", err)
	}

	// 3. Guardar en caché para futuras busquedas
	if len(result.Results) > 0 {
		if err := s.localCache.SaveWithHash(ctx, filterHash, result); err != nil {
			slog.Error("⚠️ Failed to save to cache",
				slog.String("hash", filterHash),
				slog.String("error", err.Error()))
		} else {
			slog.Info("💾 Saved to cache",
				slog.String("hash", filterHash),
				slog.Int("items_count", len(result.Results)))
		}
	} else {
		slog.Info("⚠️ Empty result, not caching", slog.String("hash", filterHash))
	}

	return result, nil

}

type ItemEvent struct {
	Action string `json:"action"` // "create", "update", "delete"
	ItemID string `json:"item_id"`
}

func (s *SearchServiceImpl) InitConsumer(ctx context.Context) {
	// Iniciar Go routine para el consumer
	slog.Info("🐰 Starting RabbitMQ consumer...")

	if err := s.consumer.Consume(ctx, s.handleMessage); err != nil {
		slog.Error("❌ Error in RabbitMQ consumer: %v", err)
	}
	slog.Info("🐰 RabbitMQ consumer stopped.")
}

// handleMessage procesa los mensajes recibidos de RabbitMQ
func (s *SearchServiceImpl) handleMessage(ctx context.Context, message ItemEvent) error {
	slog.Info("📨 Processing message",
		slog.String("action", message.Action),
		slog.String("item_id", message.ItemID),
	)

	switch message.Action {
	case "create":

		item, err := GetItemByID(message.ItemID)

		if err != nil {
			slog.Error("❌ Error getting item details",
				slog.String("item_id", message.ItemID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error getting item details: %w", err)
		}

		if _, err := s.repo.Create(ctx, item); err != nil {
			slog.Error("❌ Error indexing item in search",
				slog.String("item_id", message.ItemID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error indexing item in search: %w", err)
		}

		slog.Info("🔍 Item indexed in search engine", slog.String("item_id", message.ItemID))

	case "update":
		item, err := GetItemByID(message.ItemID)

		if err != nil {
			slog.Error("❌ Error getting item details",
				slog.String("item_id", message.ItemID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error getting item details: %w", err)
		}
		if _, err := s.repo.Update(ctx, message.ItemID, item); err != nil {
			slog.Error("❌ Error updating item in search",
				slog.String("item_id", message.ItemID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error updating item in search: %w", err)
		}

		slog.Info("✏️ Item updated", slog.String("item_id", message.ItemID))

	case "delete":
		slog.Info("🗑️ Item deleted", slog.String("item_id", message.ItemID))

		if err := s.repo.Delete(ctx, message.ItemID); err != nil {
			slog.Error("❌ Error deleting item from search",
				slog.String("item_id", message.ItemID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error deleting item from search: %w", err)
		}

		slog.Info("🗑️ Item deleted", slog.String("item_id", message.ItemID))

	default:
		slog.Info("⚠️ Unknown action", slog.String("action", message.Action))
	}

	return nil
}

func GetItemByID(id string) (domain.Item, error) {
	response, err := http.Get(fmt.Sprintf("http://products-api:8080/items/%s", id))
	if err != nil {
		return domain.Item{}, fmt.Errorf("error getting item: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return domain.Item{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error reading response: %w", err)
	}

	//Define un struct wrapper que coincida con el JSON ("item": ...)
	var responseWrapper struct {
		Item domain.Item `json:"item"`
	}

	// 2. Deserializa el JSON en el wrapper
	if err := json.Unmarshal(bytes, &responseWrapper); err != nil {
		return domain.Item{}, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return responseWrapper.Item, nil

}
