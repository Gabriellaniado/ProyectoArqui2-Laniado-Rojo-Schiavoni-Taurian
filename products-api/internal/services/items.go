package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"products-api/internal/domain"
	"strings"
)

type ItemsService interface {

	// Create valida y crea un nuevo item
	Create(ctx context.Context, item domain.Item) (domain.Item, error)

	// GetByID obtiene un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)

	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)

	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
}

// ItemsRepository define las operaciones de datos para Items
// Patrón Repository: abstrae el acceso a datos del resto de la aplicación
type ItemsRepository interface {

	// Create inserta un nuevo item en DB
	Create(ctx context.Context, item domain.Item) (domain.Item, error)

	// GetByID busca un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)

	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)

	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
} // ItemsServiceImpl implementa ItemsService

type ItemsPublisher interface { //genera mensajes para rabbit
	Publish(ctx context.Context, action string, itemID string) error
}

type ItemsServiceImpl struct {
	repository       ItemsRepository // Inyección de dependencia
	localCache       ItemsRepository // Inyección de dependencia
	distributedCache ItemsRepository // Inyección de dependencia
	publisher        ItemsPublisher
}

// NewItemsService crea una nueva instancia del service
// Pattern: Dependency Injection - recibe dependencies como parámetros
func NewItemsService(repository ItemsRepository, localCache ItemsRepository, distributedCache ItemsRepository, publisher ItemsPublisher) ItemsServiceImpl {
	return ItemsServiceImpl{
		repository:       repository,
		localCache:       localCache,
		distributedCache: distributedCache,
		publisher:        publisher,
	}
}

// Create valida y crea un nuevo item
// Consigna 1: Validar name no vacío y price >= 0
func (s *ItemsServiceImpl) Create(ctx context.Context, item domain.Item) (domain.Item, error) {

	if item.Name == "" || item.Category == "" || item.Description == "" || item.Price == 0 || item.Stock == 0 {
		return domain.Item{}, fmt.Errorf("error, all fields need to be filled")
	}
	if item.Price <= 0 {
		return domain.Item{}, fmt.Errorf("error, the price cannot be negative")
	}

	if item.Stock <= 0 {
		return domain.Item{}, fmt.Errorf("error, the stock cannot be negative")
	}

	created, err := s.repository.Create(ctx, item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error creating item in repository: %w", err)
	}

	if err := s.publisher.Publish(ctx, "create", created.ID); err != nil {
		return domain.Item{}, fmt.Errorf("error publishing item creation: %w", err)
	}

	_, err = s.distributedCache.Create(ctx, created)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error creating item in distributed cache: %w", err)
	}
	_, err = s.localCache.Create(ctx, created)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error creating item in local cache: %w", err)
	}

	return created, nil
}

// GetByID obtiene un item por su ID
// Consigna 2: Validar formato de ID antes de consultar DB
func (s *ItemsServiceImpl) GetByID(ctx context.Context, id string) (domain.Item, error) {
	item, err := s.localCache.GetByID(ctx, id)
	if err != nil {
		item, err := s.distributedCache.GetByID(ctx, id)

		if err != nil {
			item, err := s.repository.GetByID(ctx, id)
			if err != nil {
				return domain.Item{}, fmt.Errorf("error getting item from repository: %w", err)
			}

			_, err = s.localCache.Create(ctx, item)
			if err != nil {
				return domain.Item{}, fmt.Errorf("error creating item in local cache: %w", err)
			}
			_, err = s.distributedCache.Create(ctx, item)
			if err != nil {
				return domain.Item{}, fmt.Errorf("error creating item in distributed cache: %w", err)
			}

			return item, nil
		}

		s.localCache.Create(ctx, item)
		return item, nil
	}
	return item, nil
}

// Update actualiza un item existente
// Consigna 3: Validar campos antes de actualizar
func (s *ItemsServiceImpl) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {

	if item.Name == "" || item.Category == "" || item.Description == "" || item.Price == 0 || item.Stock == 0 {
		return domain.Item{}, fmt.Errorf("error, all fields need to be filled")
	}
	if item.Price < 0 {
		return domain.Item{}, fmt.Errorf("error, the price cannot be negative")
	}

	if item.Stock < 0 {
		return domain.Item{}, fmt.Errorf("error, the stock cannot be negative")
	}

	if err := s.validateItem(item); err != nil {
		return domain.Item{}, fmt.Errorf("validation error: %w", err)
	}

	// TODO: Actualizar en DB

	updated, err := s.repository.Update(ctx, id, item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error updating item in repository: %w", err)
	}

	//publicar evento de actualización

	if err := s.publisher.Publish(ctx, "update", id); err != nil {
		return domain.Item{}, fmt.Errorf("error publishing item update: %w", err)
	}

	// TODO: Guardar en cache

	if _, err = s.distributedCache.Update(ctx, id, updated); err != nil {
		slog.Warn("⚠️ Error updating item in distributed cache", slog.String("item_id", id), slog.String("error", err.Error()))
	}
	if _, err = s.localCache.Update(ctx, id, updated); err != nil {
		slog.Warn("⚠️ Error updating item in local cache", slog.String("item_id", id), slog.String("error", err.Error()))
	}

	return updated, nil
}

// Delete elimina un item por ID
// Consigna 4: Validar ID antes de eliminar
func (s *ItemsServiceImpl) Delete(ctx context.Context, id string) error {
	// TODO: Borrar de cache
	err := s.localCache.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting item from local cache: %w", err)
	}
	err = s.distributedCache.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting item from distributed cache: %w", err)
	}

	err = s.publisher.Publish(ctx, "delete", id)
	if err != nil {
		return fmt.Errorf("error publishing item deletion: %w", err)
	}

	// TODO: Borrar de DB
	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting item from repository: %w", err)
	}

	return nil
}

// validateItem aplica reglas de negocio para validar un item
// 🎯 Función helper para reutilizar validaciones
func (s *ItemsServiceImpl) validateItem(item domain.Item) error {
	// 📝 Name es obligatorio y no puede estar vacío
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}

	// 💰 Price debe ser >= 0 (productos gratis están permitidos)
	if item.Price < 0 {
		return errors.New("price must be greater than or equal to 0")
	}

	// ✅ Todas las validaciones pasaron
	return nil
}

type ItemEvent struct {
	Action string `json:"action"` // "create", "update", "delete"
	ItemID string `json:"item_id"`
}
