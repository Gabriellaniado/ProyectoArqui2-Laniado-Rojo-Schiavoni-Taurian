package repository

import (
	"context"
	"fmt"
	"log"
	"products-api/internal/dao"
	"products-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCartRepository implementa CartRepository usando MongoDB
type MongoCartRepository struct {
	collection *mongo.Collection
}

// NewMongoCartRepository crea una nueva instancia del repositorio
func NewMongoCartRepository(ctx context.Context, uri, dbName, collectionName string) *MongoCartRepository {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	// Crear índice único por customer_id para evitar múltiples carritos por usuario
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "customer_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Warning: Could not create unique index on customer_id: %v", err)
	}

	log.Printf("✓ Connected to MongoDB - Cart Collection: %s", collectionName)
	return &MongoCartRepository{collection: collection}
}

// GetByCustomerID obtiene el carrito de un cliente
func (r *MongoCartRepository) GetByCustomerID(ctx context.Context, customerID int) (domain.Cart, error) {
	var cartDAO dao.CartDAO
	filter := bson.M{"customer_id": customerID}

	err := r.collection.FindOne(ctx, filter).Decode(&cartDAO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No existe carrito, retornar un carrito vacío
			return domain.Cart{
				CustomerID: customerID,
				Items:      []domain.CartItem{},
				Total:      0,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}, nil
		}
		return domain.Cart{}, fmt.Errorf("error finding cart: %w", err)
	}

	return r.daoToDomain(cartDAO), nil
}

// Create crea un nuevo carrito
func (r *MongoCartRepository) Create(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	cartDAO := r.domainToDAO(cart)
	cartDAO.CreatedAt = time.Now()
	cartDAO.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, cartDAO)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("error creating cart: %w", err)
	}

	cart.ID = result.InsertedID.(primitive.ObjectID).Hex()
	cart.CreatedAt = cartDAO.CreatedAt
	cart.UpdatedAt = cartDAO.UpdatedAt
	return cart, nil
}

// Update actualiza un carrito existente
func (r *MongoCartRepository) Update(ctx context.Context, customerID int, cart domain.Cart) (domain.Cart, error) {
	cartDAO := r.domainToDAO(cart)
	cartDAO.UpdatedAt = time.Now()

	filter := bson.M{"customer_id": customerID}
	update := bson.M{
		"$set": bson.M{
			"items":      cartDAO.Items,
			"total":      cartDAO.Total,
			"updated_at": cartDAO.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedDAO dao.CartDAO
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDAO)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("error updating cart: %w", err)
	}

	return r.daoToDomain(updatedDAO), nil
}

// Delete elimina un carrito (por ejemplo, después del checkout)
func (r *MongoCartRepository) Delete(ctx context.Context, customerID int) error {
	filter := bson.M{"customer_id": customerID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting cart: %w", err)
	}
	return nil
}

// Upsert crea o actualiza un carrito
func (r *MongoCartRepository) Upsert(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	cartDAO := r.domainToDAO(cart)
	now := time.Now()
	cartDAO.UpdatedAt = now

	filter := bson.M{"customer_id": cart.CustomerID}
	update := bson.M{
		"$set": bson.M{
			"items":      cartDAO.Items,
			"total":      cartDAO.Total,
			"updated_at": cartDAO.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"customer_id": cart.CustomerID,
			"created_at":  now,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var updatedDAO dao.CartDAO
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDAO)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("error upserting cart: %w", err)
	}

	return r.daoToDomain(updatedDAO), nil
}

// Conversión de DAO a Domain
func (r *MongoCartRepository) daoToDomain(cartDAO dao.CartDAO) domain.Cart {
	items := make([]domain.CartItem, len(cartDAO.Items))
	for i, item := range cartDAO.Items {
		items[i] = domain.CartItem{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		}
	}

	return domain.Cart{
		ID:         cartDAO.ID.Hex(),
		CustomerID: cartDAO.CustomerID,
		Items:      items,
		Total:      cartDAO.Total,
		CreatedAt:  cartDAO.CreatedAt,
		UpdatedAt:  cartDAO.UpdatedAt,
	}
}

// Conversión de Domain a DAO
func (r *MongoCartRepository) domainToDAO(cart domain.Cart) dao.CartDAO {
	items := make([]dao.CartItemDAO, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = dao.CartItemDAO{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		}
	}

	var objectID primitive.ObjectID
	if cart.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(cart.ID)
	}

	return dao.CartDAO{
		ID:         objectID,
		CustomerID: cart.CustomerID,
		Items:      items,
		Total:      cart.Total,
		CreatedAt:  cart.CreatedAt,
		UpdatedAt:  cart.UpdatedAt,
	}
}
