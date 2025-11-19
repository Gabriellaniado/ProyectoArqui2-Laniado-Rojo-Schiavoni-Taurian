package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"products-api/internal/dao"
	"products-api/internal/domain"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// coleccions de items en mongo
type MongoItemsRepository struct {
	col *mongo.Collection
}

// conexion a mongo y a la coleccion items
func NewMongoItemsRepository(ctx context.Context, uri, dbName, collectionName string) *MongoItemsRepository {
	opt := options.Client().ApplyURI(uri)
	opt.SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
		return nil
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
		return nil
	}

	return &MongoItemsRepository{
		col: client.Database(dbName).Collection(collectionName), // Conecta con la colección "items"
	}
}

//get by id

func (r *MongoItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	// Validar que el ID tenga formato ObjectID válido
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Item{}, errors.New("invalid ObjectID format")
	}

	// Buscar el documento por _id
	var daoItem dao.Item
	filter := bson.M{"_id": objectID}
	err = r.col.FindOne(ctx, filter).Decode(&daoItem)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Item{}, errors.New("item not found")
		}
		return domain.Item{}, err
	}

	return daoItem.ToDomain(), nil
}

// Devuelve items paginados

func (r *MongoItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	itemDAO := dao.FromDomain(item) // Convertir a DAO para manejar ObjectID y BSON
	itemDAO.ID = primitive.NewObjectID()
	now := time.Now().UTC().Truncate(time.Millisecond)
	itemDAO.CreatedAt = now
	itemDAO.UpdatedAt = now

	// Insertar en DB
	res, err := r.col.InsertOne(ctx, itemDAO)
	if err != nil {
		return domain.Item{}, err
	}

	// Obtener el ID generado por DB y convertir a string
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		itemDAO.ID = oid
	} else {
		return domain.Item{}, errors.New("failed to convert inserted ID to ObjectID")
	}

	return itemDAO.ToDomain(), nil // Convertir de vuelta a Domain para retornar
}

// Update actualiza un item existente
// Consigna 3: Update parcial + actualizar updatedAt
func (r *MongoItemsRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// convertir el id string a ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Item{}, errors.New("invalid ObjectID format")
	}

	updateFields := bson.M{
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
		"stock":       item.Stock,
		"category":    item.Category,
		"image_url":   item.ImageURL,
		"updated_at":  time.Now().UTC().Truncate(time.Millisecond), // Solo actualizar updated_at
	}

	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateFields})
	if err != nil {
		return domain.Item{}, err
	}

	// obtener el documento actualizado y devolverlo como domain.Item (con id)
	var updatedDAO dao.Item
	if err := r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedDAO); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Item{}, errors.New("item not found")
		}
		return domain.Item{}, err
	}

	domainItem := updatedDAO.ToDomain()

	return domainItem, nil
}

// Delete elimina un item por ID
// Consigna 4: Eliminar documento de DB
func (r *MongoItemsRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		return err
	}

	return nil
}

// DecrementStockAtomic decrementa stock SOLO si hay suficiente (operación atómica)
func (r *MongoItemsRepository) DecrementStockAtomic(ctx context.Context, itemID string, quantity int) (bool, error) {

	// Convertir string a ObjectID
	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		log.Printf("❌ Invalid ObjectID format: %s", itemID)
		return false, errors.New("invalid ObjectID format")
	}
	filter := bson.M{
		"_id":   objID,
		"stock": bson.M{"$gte": quantity}, // Solo si stock >= quantity
	}
	update := bson.M{
		"$inc": bson.M{"stock": -quantity}, //  Decrementar
	}

	result := r.col.FindOneAndUpdate(ctx, filter, update)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil // No había stock suficiente
		}
		return false, err
	}
	return true, nil
}

// IncrementStock incrementa stock (para rollback)
func (r *MongoItemsRepository) IncrementStock(ctx context.Context, itemID string, quantity int) error {
	objID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		log.Printf("❌ Invalid ObjectID format: %s", itemID)
		return errors.New("invalid ObjectID format")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"stock": quantity}}
	_, err = r.col.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Printf("❌ Error incrementing stock: %v", err)
		return err
	}

	log.Printf("✅ Stock incremented for item %s", itemID)
	return nil
}
