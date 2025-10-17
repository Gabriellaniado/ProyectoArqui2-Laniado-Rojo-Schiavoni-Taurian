package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"products-api/internal/dao"
	"products-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoSalesRepository maneja las operaciones de ventas en MongoDB
type MongoSalesRepository struct {
	col *mongo.Collection
}

// NewMongoSalesRepository crea una nueva instancia del repository y se conecta a MongoDB
func NewMongoSalesRepository(ctx context.Context, uri, dbName, collectionName string) *MongoSalesRepository {
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

	return &MongoSalesRepository{
		col: client.Database(dbName).Collection(collectionName),
	}
}

// GetByID obtiene una venta por su ID de MongoDB
func (r *MongoSalesRepository) GetByID(ctx context.Context, id string) (domain.Sales, error) {
	// Validar que el ID tenga formato ObjectID válido
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Sales{}, errors.New("invalid ObjectID format")
	}

	// Buscar el documento por _id
	var daoSale dao.Sales
	filter := bson.M{"_id": objectID}
	err = r.col.FindOne(ctx, filter).Decode(&daoSale)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Sales{}, errors.New("sale not found")
		}
		return domain.Sales{}, err
	}

	return daoSale.ToDomain(), nil
}

// GetBySaleID obtiene una venta por su SaleID (ID de negocio)
func (r *MongoSalesRepository) GetBySaleID(ctx context.Context, saleID string) (domain.Sales, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var daoSale dao.Sales
	filter := bson.M{"sale_id": saleID}
	err := r.col.FindOne(ctx, filter).Decode(&daoSale)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Sales{}, errors.New("sale not found")
		}
		return domain.Sales{}, err
	}

	return daoSale.ToDomain(), nil
}

// En internal/repository/sales_mongo.go

// GetByCustomerID obtiene todas las ventas de un cliente específico
func (r *MongoSalesRepository) GetByCustomerID(ctx context.Context, customerID string) ([]domain.Sales, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Filtrar por customer_id
	filter := bson.M{"customer_id": customerID}

	// Buscar todos los documentos que coincidan
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return []domain.Sales{}, err
	}
	defer cur.Close(ctx)

	// Decodificar resultados
	var daoSales []dao.Sales
	if err := cur.All(ctx, &daoSales); err != nil {
		return []domain.Sales{}, err
	}

	// Si no hay ventas, retornar slice vacío (no error)
	if len(daoSales) == 0 {
		return []domain.Sales{}, nil
	}

	// Convertir de DAO a Domain
	domainSales := make([]domain.Sales, len(daoSales))
	for i, daoSale := range daoSales {
		domainSales[i] = daoSale.ToDomain()
	}

	return domainSales, nil
}

// Create crea una nueva venta en la base de datos
func (r *MongoSalesRepository) Create(ctx context.Context, sale domain.Sales) (domain.Sales, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	saleDAO := dao.FromDomainSales(sale)
	saleDAO.ID = primitive.NewObjectID()

	// Si SaleDate no está seteada, usar la fecha actual
	saleDAO.SaleDate = time.Now().UTC()

	// Insertar en DB
	res, err := r.col.InsertOne(ctx, saleDAO)
	if err != nil {
		return domain.Sales{}, err
	}

	// Obtener el ID generado por DB
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		saleDAO.ID = oid
	} else {
		return domain.Sales{}, errors.New("failed to convert inserted ID to ObjectID")
	}

	return saleDAO.ToDomain(), nil
}

// Update actualiza una venta existente
func (r *MongoSalesRepository) Update(ctx context.Context, id string, sale domain.Sales) (domain.Sales, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validar ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Sales{}, errors.New("invalid ObjectID format")
	}

	// Convertir a DAO
	daoSale := dao.FromDomainSales(sale)

	// Actualizar en DB
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": daoSale}

	result, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Sales{}, err
	}

	if result.MatchedCount == 0 {
		return domain.Sales{}, errors.New("sale not found")
	}

	// Retornar la venta actualizada
	daoSale.ID = objectID
	return daoSale.ToDomain(), nil
}

// Delete elimina una venta por ID
func (r *MongoSalesRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ObjectID format")
	}

	result, err := r.col.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("sale not found")
	}

	return nil
}
