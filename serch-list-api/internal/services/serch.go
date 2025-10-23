package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"serch-list-api/internal/domain"
)

type SerchRepository interface {
	List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error)
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)
	Delete(ctx context.Context, id string) error
}

type ItemsConsumer interface { //consume mensajes de rabbit
	Consume(ctx context.Context, handler func(ctx context.Context, message ItemEvent) error) error
}

type SerchServiceImpl struct {
	repo     SerchRepository
	consumer ItemsConsumer
}

func NewSerchService(repo SerchRepository, consumer ItemsConsumer) *SerchServiceImpl {
	return &SerchServiceImpl{
		repo:     repo,
		consumer: consumer,
	}
}

func (s *SerchServiceImpl) List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error) {
	return s.repo.List(ctx, filters)
}

type ItemEvent struct {
	Action string `json:"action"` // "create", "update", "delete"
	ItemID string `json:"item_id"`
}

func (s *SerchServiceImpl) InitConsumer(ctx context.Context) {
	// Iniciar Go routine para el consumer
	slog.Info("🐰 Starting RabbitMQ consumer...")

	if err := s.consumer.Consume(ctx, s.handleMessage); err != nil {
		slog.Error("❌ Error in RabbitMQ consumer: %v", err)
	}
	slog.Info("🐰 RabbitMQ consumer stopped.")
}

// handleMessage procesa los mensajes recibidos de RabbitMQ
func (s *SerchServiceImpl) handleMessage(ctx context.Context, message ItemEvent) error {
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
		if _, err := s.repo.Update(ctx, message.ItemID, domain.Item{}); err != nil {
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
	response, err := http.Get(fmt.Sprintf("http://localhost:8080/items/%s", id))
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
