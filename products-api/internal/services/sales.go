package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"products-api/internal/domain"
	"strings"
)

// SalesRepository define las operaciones de datos para Sales
type SalesRepository interface {
	Create(ctx context.Context, sale domain.Sales) (domain.Sales, error)
	GetByID(ctx context.Context, id string) (domain.Sales, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]domain.Sales, error)
	Update(ctx context.Context, id string, sale domain.Sales) (domain.Sales, error)
	Delete(ctx context.Context, id string) error
}

// SalesServiceImpl implementa SalesService
type SalesServiceImpl struct {
<<<<<<< HEAD
<<<<<<< HEAD
	repository   SalesRepository
	localCache   SalesRepository
	itemsService ItemsService
=======
	repository  SalesRepository // Inyección de dependencia
	cache       SalesRepository // Inyección de dependencia
	usersAPIURL string
>>>>>>> main
=======
	repository   SalesRepository
	localCache   SalesRepository
	itemsService ItemsService
>>>>>>> main
}

// NewSalesService crea una nueva instancia del service
func NewSalesService(repository SalesRepository, cache SalesRepository, itemsService ItemsService) SalesServiceImpl {
	return SalesServiceImpl{
<<<<<<< HEAD
<<<<<<< HEAD
		repository:   repository,
		localCache:   cache,
		itemsService: itemsService,
=======
		repository:  repository,
		cache:       cache,
		usersAPIURL: "http://localhost:8082",
>>>>>>> main
=======
		repository:   repository,
		localCache:   cache,
		itemsService: itemsService,
>>>>>>> main
	}
}

// Create valida y crea una nueva venta, decrementando el stock del item
func (s *SalesServiceImpl) Create(ctx context.Context, sale domain.BodySales) (domain.Sales, error) {
	// Validar la venta antes de crearla
	if err := s.validateSale(sale); err != nil {
		return domain.Sales{}, fmt.Errorf("validation error: %w", err)
	}

	// Obtener el item para validar stock y calcular precio
	item, err := s.itemsService.GetByID(ctx, sale.ItemID)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error getting item: %w", err)
	}

	// Validar que el item exista
	if item.ID == "" {
		return domain.Sales{}, errors.New("item does not exist")
	}

	// Validar que haya stock suficiente
	if item.Stock < sale.Quantity {
		return domain.Sales{}, fmt.Errorf("insufficient stock: requested %d, available %d", sale.Quantity, item.Stock)
	}

	newSale := domain.Sales{
		ItemID:     sale.ItemID,
		Quantity:   sale.Quantity,
		TotalPrice: 0, // Se calculará abajo
		CustomerID: sale.CustomerID,
	}

	// Calcular el precio total de la venta
	newSale.TotalPrice = item.Price * float64(newSale.Quantity)

	// Decrementar el stock del item
	item.Stock -= sale.Quantity
	_, err = s.itemsService.Update(ctx, sale.ItemID, item)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error updating item stock: %w", err)
	}

	item, err = s.itemsService.Update(ctx, sale.ItemID, item)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error updating item stock: %w", err)
	}

	// Crear la venta en el repository
	created, err := s.repository.Create(ctx, newSale)
	if err != nil {
		//  Si falla la creación de la venta, intentar revertir el stock
		// (En un sistema real, esto debería manejarse con transacciones)
		item.Stock += sale.Quantity
		if revertErr := s.revertStock(ctx, sale.ItemID, item); revertErr != nil {
			log.Printf(" Failed to revert stock after sale creation error: %v", revertErr)
		}
		return domain.Sales{}, fmt.Errorf("error creating sale in repository: %w", err)
	}

	// Guardar en cache
	_, err = s.localCache.Create(ctx, created)
	if err != nil {
		log.Printf(" Error creating sale in cache: %v", err)
	}

	return created, nil
}

// revertStock intenta revertir el stock en caso de error
func (s *SalesServiceImpl) revertStock(ctx context.Context, itemID string, item domain.Item) error {
	_, err := s.itemsService.Update(ctx, itemID, item)
	return err
}

// GetByID obtiene una venta por su ID
func (s *SalesServiceImpl) GetByID(ctx context.Context, id string) (domain.Sales, error) {
	// Intentar obtener de cache primero
	sale, err := s.localCache.GetByID(ctx, id)
	if err != nil {
		// Si no está en cache, buscar en repository
		sale, err = s.repository.GetByID(ctx, id)
		log.Println(" Cache MISS - Sale fetched from repository:", sale.ID)
		if err != nil {
			return domain.Sales{}, fmt.Errorf("error getting sale from repository: %w", err)
		}

		// Guardar en cache para futuras consultas
		_, err = s.localCache.Create(ctx, sale)
		if err != nil {
			log.Printf(" Error creating sale in cache: %v", err)
		}

		return sale, nil
	}
	log.Println(" Cache HIT - Sale fetched from cache:", sale.ID)
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
func (s *SalesServiceImpl) Update(ctx context.Context, id string, sale domain.BodySales) (domain.Sales, error) {
	if err := s.validateSale(sale); err != nil {
		return domain.Sales{}, fmt.Errorf("validation error: %w", err)
	}

	// Obtener la venta original para calcular la diferencia de stock
	originalSale, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error getting original sale: %w", err)
	}

	// Obtener el item para validar stock
	item, err := s.itemsService.GetByID(ctx, sale.ItemID)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error getting item: %w", err)
	}

	// Calcular la diferencia de cantidad
	quantityDiff := sale.Quantity - originalSale.Quantity

	// Si se incrementó la cantidad, validar stock disponible
	if quantityDiff > 0 {
		if item.Stock < quantityDiff {
			return domain.Sales{}, fmt.Errorf("insufficient stock: need %d more, available %d", quantityDiff, item.Stock)
		}
		item.Stock -= quantityDiff
	} else if quantityDiff < 0 {
		// Si se decrementó la cantidad, devolver stock
		item.Stock += -quantityDiff
	}

	item, err = s.itemsService.Update(ctx, sale.ItemID, item)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error updating item stock: %w", err)
	}

	newSale := domain.Sales{
		ItemID:     sale.ItemID,
		Quantity:   sale.Quantity,
		TotalPrice: 0, // Se recalculará abajo
		CustomerID: sale.CustomerID,
	}

	// Recalcular el precio total
	newSale.TotalPrice = item.Price * float64(newSale.Quantity)

	// Actualizar el stock si hubo cambios
	if quantityDiff != 0 {
		_, err = s.itemsService.Update(ctx, sale.ItemID, item)
		if err != nil {
			return domain.Sales{}, fmt.Errorf("error updating item stock: %w", err)
		}
	}

	//  Actualizar la venta
	updated, err := s.repository.Update(ctx, id, newSale)
	if err != nil {
		return domain.Sales{}, fmt.Errorf("error updating sale in repository: %w", err)
	}

	// Actualizar en cache
	_, err = s.localCache.Update(ctx, id, updated)
	if err != nil {
		log.Printf("Error updating sale in cache: %v", err)
	}

	return updated, nil
}

// Delete elimina una venta por ID y restaura el stock
func (s *SalesServiceImpl) Delete(ctx context.Context, id string) error {
	// Obtener la venta para saber cuánto stock restaurar
	sale, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting sale: %w", err)
	}

	// Obtener el item para restaurar el stock
	item, err := s.itemsService.GetByID(ctx, sale.ItemID)
	if err != nil {
		return fmt.Errorf("error getting item: %w", err)
	}

	// Restaurar el stock
	item.Stock += sale.Quantity
	_, err = s.itemsService.Update(ctx, sale.ItemID, item)
	if err != nil {
		return fmt.Errorf("error restoring item stock: %w", err)
	}

	// Borrar de cache
	err = s.localCache.Delete(ctx, id)
	if err != nil {
		log.Printf(" Error deleting sale from cache: %v", err)
	}

	// Borrar de DB
	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf(" error deleting sale from repository: %w", err)
	}

	return nil
}

// validateSale aplica reglas de negocio para validar una venta
func (s *SalesServiceImpl) validateSale(sale domain.BodySales) error {
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

	// TotalPrice no puede ser negativo (si viene en el request)

	return nil
}
