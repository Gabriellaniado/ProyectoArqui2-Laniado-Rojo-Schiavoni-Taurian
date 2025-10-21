package services

import (
	"context"
	"errors"
	"fmt"
	"products-api/internal/domain"
	"strings"
)

// SalesRepository define las operaciones de datos para Sales
// Patrón Repository: abstrae el acceso a datos del resto de la aplicación
type SalesRepository interface {

	// Create inserta una nueva venta en DB
	Create(ctx context.Context, sale domain.Sales) (domain.Sales, error)

	// GetByID busca una venta por su ID
	GetByID(ctx context.Context, id string) (domain.Sales, error)

	// GetByCustomerID busca todas las ventas de un cliente
	GetByCustomerID(ctx context.Context, customerID string) ([]domain.Sales, error)

	// Update actualiza una venta existente
	Update(ctx context.Context, id string, sale domain.Sales) (domain.Sales, error)

	// Delete elimina una venta por ID
	Delete(ctx context.Context, id string) error
}

// SalesServiceImpl implementa SalesService
type SalesServiceImpl struct {
	repository SalesRepository // Inyección de dependencia
	cache      SalesRepository // Inyección de dependencia
}

// NewSalesService crea una nueva instancia del service
// Pattern: Dependency Injection - recibe dependencies como parámetros
func NewSalesService(repository SalesRepository, cache SalesRepository) SalesServiceImpl {
	return SalesServiceImpl{
		repository: repository,
		cache:      cache,
	}
}

// Create valida y crea una nueva venta
func (s *SalesServiceImpl) Create(ctx context.Context, sale domain.Sales) (domain.Sales, error) {
	// Validar la venta antes de crearla
	if err := s.validateSale(sale); err != nil {
		return domain.Sales{}, fmt.Errorf("validation error: %w", err)
	}

	created, err := s.repository.Create(ctx, sale)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error creating sale in repository: %w", err)
	}

	_, err = s.cache.Create(ctx, created)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error creating sale in cache: %w", err)
	}

	return created, nil
}

// GetByID obtiene una venta por su ID
func (s *SalesServiceImpl) GetByID(ctx context.Context, id string) (domain.Sales, error) {
	// Intentar obtener de cache primero
	sale, err := s.cache.GetByID(ctx, id)
	if err != nil {
		// Si no está en cache, buscar en repository
		sale, err := s.repository.GetByID(ctx, id)
		if err != nil {
			return domain.Sales{}, fmt.Errorf("error getting sale from repository: %w", err)
		}

		// Guardar en cache para futuras consultas
		_, err = s.cache.Create(ctx, sale)
		if err != nil {
			return domain.Sales{}, fmt.Errorf("error creating sale in cache: %w", err)
		}

		return sale, nil
	}
	return sale, nil
}

// GetByCustomerID obtiene todas las ventas de un cliente específico
func (s *SalesServiceImpl) GetByCustomerID(ctx context.Context, customerID string) ([]domain.Sales, error) {
	sales, err := s.repository.GetByCustomerID(ctx, customerID)
	if err != nil {
		return []domain.Sales{}, fmt.Errorf("error getting sales by customer_id from repository: %w", err)
	}
	return sales, nil
}

// Update actualiza una venta existente
func (s *SalesServiceImpl) Update(ctx context.Context, id string, sale domain.Sales) (domain.Sales, error) {
	if err := s.validateSale(sale); err != nil {
		return domain.Sales{}, fmt.Errorf("validation error: %w", err)
	}

	updated, err := s.repository.Update(ctx, id, sale)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error updating sale in repository: %w", err)
	}

	_, err = s.cache.Update(ctx, id, sale)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error updating sale in cache: %w", err)
	}

	return updated, nil
}

// Delete elimina una venta por ID
func (s *SalesServiceImpl) Delete(ctx context.Context, id string) error {
	// Borrar de cache
	err := s.cache.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting sale from cache: %w", err)
	}

	// Borrar de DB
	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting sale from repository: %w", err)
	}

	return nil
}

// validateSale aplica reglas de negocio para validar una venta
func (s *SalesServiceImpl) validateSale(sale domain.Sales) error {
	// SaleID es obligatorio

	// ItemID es obligatorio
	if strings.TrimSpace(sale.ItemID) == "" {
		return errors.New("item_id is required and cannot be empty")
	}

	// CustomerID es obligatorio
	if strings.TrimSpace(sale.CustomerID) == "" {
		return errors.New("customer_id is required and cannot be empty")
	}

	// Quantity debe ser mayor a 0
	if sale.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// TotalPrice debe ser >= 0
	if sale.TotalPrice < 0 {
		return errors.New("total_price must be greater than or equal to 0")
	}

	return nil
}
