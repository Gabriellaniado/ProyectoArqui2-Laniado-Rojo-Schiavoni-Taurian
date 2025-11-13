package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"products-api/internal/domain"
)

// CartRepository define las operaciones de datos para Cart
type CartRepository interface {
	GetByCustomerID(ctx context.Context, customerID int) (domain.Cart, error)
	Create(ctx context.Context, cart domain.Cart) (domain.Cart, error)
	Update(ctx context.Context, customerID int, cart domain.Cart) (domain.Cart, error)
	Delete(ctx context.Context, customerID int) error
	Upsert(ctx context.Context, cart domain.Cart) (domain.Cart, error)
}

// CartServiceImpl implementa CartService
type CartServiceImpl struct {
	repository   CartRepository
	localCache   CartRepository
	itemsService ItemsService
	salesService *SalesServiceImpl
}

// NewCartService crea una nueva instancia del service
func NewCartService(repository CartRepository, cache CartRepository, itemsService ItemsService, salesService *SalesServiceImpl) *CartServiceImpl {
	return &CartServiceImpl{
		repository:   repository,
		localCache:   cache,
		itemsService: itemsService,
		salesService: salesService,
	}
}

// GetCart obtiene el carrito de un cliente con información enriquecida
func (s *CartServiceImpl) GetCart(ctx context.Context, customerID int) (domain.CartResponse, error) {
	// Intentar obtener del cache primero
	cart, err := s.localCache.GetByCustomerID(ctx, customerID)
	if err != nil {
		// Si no está en cache, buscar en repository
		cart, err = s.repository.GetByCustomerID(ctx, customerID)
		if err != nil {
			return domain.CartResponse{}, fmt.Errorf("error getting cart from repository: %w", err)
		}
		log.Printf("🔍 Cache MISS - Cart fetched from repository for customer: %d", customerID)

		// Guardar en cache
		_, _ = s.localCache.Create(ctx, cart)
	} else {
		log.Printf("✅ Cache HIT - Cart fetched from cache for customer: %d", customerID)
	}

	// Enriquecer el carrito con información de los productos
	return s.enrichCart(ctx, cart)
}

// AddItem agrega un producto al carrito o incrementa su cantidad
func (s *CartServiceImpl) AddItem(ctx context.Context, customerID int, req domain.AddItemRequest) (domain.CartResponse, error) {
	// Validar que el producto exista y tenga stock
	item, err := s.itemsService.GetByID(ctx, req.ItemID)
	if err != nil {
		return domain.CartResponse{}, fmt.Errorf("error getting item: %w", err)
	}

	if item.ID == "" {
		return domain.CartResponse{}, errors.New("item does not exist")
	}

	if item.Stock < req.Quantity {
		return domain.CartResponse{}, fmt.Errorf("insufficient stock: requested %d, available %d", req.Quantity, item.Stock)
	}

	// Obtener carrito actual
	cart, err := s.repository.GetByCustomerID(ctx, customerID)
	if err != nil {
		// Si no existe, crear uno nuevo
		cart = domain.Cart{
			CustomerID: customerID,
			Items:      []domain.CartItem{},
			Total:      0,
		}
	}

	// Buscar si el item ya está en el carrito
	found := false
	for i, cartItem := range cart.Items {
		if cartItem.ItemID == req.ItemID {
			// Validar que la cantidad total no exceda el stock
			newQuantity := cartItem.Quantity + req.Quantity
			if newQuantity > item.Stock {
				return domain.CartResponse{}, fmt.Errorf("insufficient stock: total quantity %d exceeds available %d", newQuantity, item.Stock)
			}
			cart.Items[i].Quantity = newQuantity
			found = true
			break
		}
	}

	// Si no está en el carrito, agregarlo
	if !found {
		cart.Items = append(cart.Items, domain.CartItem{
			ItemID:   req.ItemID,
			Quantity: req.Quantity,
		})
	}

	// Actualizar en la base de datos (upsert)
	cart, err = s.repository.Upsert(ctx, cart)
	if err != nil {
		return domain.CartResponse{}, fmt.Errorf("error updating cart: %w", err)
	}

	// Actualizar cache
	_, _ = s.localCache.Upsert(ctx, cart)

	log.Printf("Item added to cart - Customer: %d, Item: %s, Quantity: %d", customerID, req.ItemID, req.Quantity)

	// Retornar carrito enriquecido
	return s.enrichCart(ctx, cart)
}

// UpdateItem actualiza la cantidad de un producto en el carrito
func (s *CartServiceImpl) UpdateItemCart(ctx context.Context, customerID int, itemID string, req domain.UpdateItemRequest) (domain.CartResponse, error) {
	// Obtener carrito actual
	cart, err := s.repository.GetByCustomerID(ctx, customerID)
	if err != nil {
		return domain.CartResponse{}, fmt.Errorf("cart not found: %w", err)
	}

	// Si la cantidad es 0, eliminar el item
	if req.Quantity == 0 {
		return s.RemoveItem(ctx, customerID, itemID)
	}

	// Validar stock disponible
	item, err := s.itemsService.GetByID(ctx, itemID)
	if err != nil {
		return domain.CartResponse{}, fmt.Errorf("error getting item: %w", err)
	}

	if item.Stock < req.Quantity {
		return domain.CartResponse{}, fmt.Errorf("insufficient stock: requested %d, available %d", req.Quantity, item.Stock)
	}

	// Buscar el item en el carrito y actualizar
	found := false
	for i, cartItem := range cart.Items {
		if cartItem.ItemID == itemID {
			cart.Items[i].Quantity = req.Quantity
			found = true
			break
		}
	}

	if !found {
		return domain.CartResponse{}, errors.New("item not found in cart")
	}

	// Actualizar en la base de datos
	cart, err = s.repository.Update(ctx, customerID, cart)
	if err != nil {
		return domain.CartResponse{}, fmt.Errorf("error updating cart: %w", err)
	}

	// Actualizar cache
	_, _ = s.localCache.Update(ctx, customerID, cart)

	log.Printf("✏️ Item updated in cart - Customer: %d, Item: %s, New Quantity: %d", customerID, itemID, req.Quantity)

	return s.enrichCart(ctx, cart)
}

// RemoveItem elimina un producto del carrito
func (s *CartServiceImpl) RemoveItem(ctx context.Context, customerID int, itemID string) (domain.CartResponse, error) {
	// Obtener carrito actual
	cart, err := s.repository.GetByCustomerID(ctx, customerID)
	if err != nil {
		return domain.CartResponse{}, fmt.Errorf("cart not found: %w", err)
	}

	// Filtrar el item a eliminar
	newItems := []domain.CartItem{}
	found := false
	for _, cartItem := range cart.Items {
		if cartItem.ItemID != itemID {
			newItems = append(newItems, cartItem)
		} else {
			found = true
		}
	}

	if !found {
		return domain.CartResponse{}, errors.New("item not found in cart")
	}

	cart.Items = newItems

	// Si el carrito quedó vacío, podríamos eliminarlo completamente
	// Pero es mejor dejarlo vacío para mantener la referencia
	cart, err = s.repository.Update(ctx, customerID, cart)
	if err != nil {
		return domain.CartResponse{}, fmt.Errorf("error updating cart: %w", err)
	}

	// Actualizar cache
	_, _ = s.localCache.Update(ctx, customerID, cart)

	log.Printf("🗑️ Item removed from cart - Customer: %d, Item: %s", customerID, itemID)

	return s.enrichCart(ctx, cart)
}

// ClearCart vacía el carrito de un cliente
func (s *CartServiceImpl) ClearCart(ctx context.Context, customerID int) error {
	cart := domain.Cart{
		CustomerID: customerID,
		Items:      []domain.CartItem{},
		Total:      0,
	}

	_, err := s.repository.Upsert(ctx, cart)
	if err != nil {
		return fmt.Errorf("error clearing cart: %w", err)
	}

	// Limpiar cache
	_ = s.localCache.Delete(ctx, customerID)

	log.Printf("🧹 Cart cleared for customer: %d", customerID)
	return nil
}

// Checkout procesa la compra del carrito
func (s *CartServiceImpl) Checkout(ctx context.Context, customerID int) ([]domain.Sales, error) {
	// Obtener carrito
	cart, err := s.repository.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validar stock de todos los items antes de procesar
	for _, cartItem := range cart.Items {
		item, err := s.itemsService.GetByID(ctx, cartItem.ItemID)
		if err != nil {
			return nil, fmt.Errorf("error validating item %s: %w", cartItem.ItemID, err)
		}

		if item.Stock < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for item %s: requested %d, available %d", item.Name, cartItem.Quantity, item.Stock)
		}
	}

	// Crear una venta por cada item del carrito
	sales := []domain.Sales{}

	for _, cartItem := range cart.Items {
		saleData := domain.BodySales{
			ItemID:     cartItem.ItemID,
			Quantity:   cartItem.Quantity,
			CustomerID: fmt.Sprintf("%d", customerID),
		}

		// Crear la venta usando el salesService
		sale, err := s.salesService.Create(ctx, saleData)
		if err != nil {
			// Si falla alguna venta, intentar revertir las ventas ya creadas
			log.Printf("❌ Error creating sale during checkout: %v", err)

			// Revertir ventas ya creadas
			for _, createdSale := range sales {
				if deleteErr := s.salesService.Delete(ctx, createdSale.ID); deleteErr != nil {
					log.Printf("⚠️ Error reverting sale %s: %v", createdSale.ID, deleteErr)
				}
			}

			return nil, fmt.Errorf("error creating sale for item %s: %w", cartItem.ItemID, err)
		}

		sales = append(sales, sale)
		log.Printf("✅ Sale created - ID: %s, Item: %s, Quantity: %d", sale.ID, cartItem.ItemID, cartItem.Quantity)
	}

	// Vaciar el carrito después del checkout exitoso
	err = s.ClearCart(ctx, customerID)
	if err != nil {
		log.Printf("⚠️ Warning: Could not clear cart after checkout: %v", err)
	}

	log.Printf("🎉 Checkout completed for customer: %d, total items: %d, total sales: %d", customerID, len(cart.Items), len(sales))
	return sales, nil
}

// enrichCart enriquece el carrito con información completa de los productos
func (s *CartServiceImpl) enrichCart(ctx context.Context, cart domain.Cart) (domain.CartResponse, error) {
	itemsWithDetails := []domain.CartItemWithDetails{}
	totalItems := 0

	for _, cartItem := range cart.Items {
		// Obtener información completa del producto
		item, err := s.itemsService.GetByID(ctx, cartItem.ItemID)
		if err != nil {
			log.Printf("⚠️ Warning: Could not get item details for %s: %v", cartItem.ItemID, err)
			// Continuar con los demás items
			continue
		}

		itemWithDetails := domain.CartItemWithDetails{
			ItemID:      cartItem.ItemID,                         // Traigo ID del cart
			Name:        item.Name,                               // Traigo nombre del producto
			Description: item.Description,                        // Traigo descripción del producto
			ImageURL:    item.ImageURL,                           // Traigo imagen del producto
			Price:       item.Price,                              // Usar el precio actual del producto
			Quantity:    cartItem.Quantity,                       // Traigo cantidad del cart
			Subtotal:    item.Price * float64(cartItem.Quantity), // Calculo subtotal con precio actual
			Stock:       item.Stock,                              // Traigo stock actual del producto
		}

		itemsWithDetails = append(itemsWithDetails, itemWithDetails)
		totalItems += cartItem.Quantity
	}

	// Recalcular el total basado en los precios actuales
	total := 0.0
	for _, item := range itemsWithDetails {
		total += item.Subtotal
	}

	return domain.CartResponse{
		ID:         cart.ID,
		CustomerID: cart.CustomerID,
		Items:      itemsWithDetails,
		Total:      total,
		ItemCount:  totalItems,
		CreatedAt:  cart.CreatedAt,
		UpdatedAt:  cart.UpdatedAt,
	}, nil
}
