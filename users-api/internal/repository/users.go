package repository

import (
	"context"
	"errors"
	"users-api/internal/dao"
	"users-api/internal/domain"

	"gorm.io/gorm"
)

// UsersRepository define las operaciones de datos para Users
// Patrón Repository: abstrae el acceso a datos del resto de la aplicación
type UsersRepository interface {
	// List retorna todos los usuarios de la base de datos
	List(ctx context.Context) ([]domain.UserDto, error)

	// Create inserta un nuevo usuario en MySQL
	Create(ctx context.Context, user domain.UserDto) (domain.UserDto, error)

	// GetByID busca un usuario por su ID
	GetByID(ctx context.Context, id int) (domain.UserDto, error)

	// GetByEmail busca un usuario por su email
	GetByEmail(ctx context.Context, email string) (domain.UserDto, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, id int, user domain.UserDto) (domain.UserDto, error)

	// Delete elimina un usuario por ID
	Delete(ctx context.Context, id int) error
}

// MySQLUsersRepository implementa UsersRepository usando MySQL con GORM
type MySQLUsersRepository struct {
	db *gorm.DB // Referencia a la base de datos MySQL
}

// NewMySQLUsersRepository crea una nueva instancia del repository
// Recibe una referencia a la base de datos MySQL
func NewMySQLUsersRepository(db *gorm.DB) UsersRepository {
	return &MySQLUsersRepository{db: db}
}

// List obtiene todos los usuarios de MySQL
func (r *MySQLUsersRepository) List(ctx context.Context) ([]domain.UserDto, error) {
	var daoUsers []dao.UserModel

	// Usar GORM para obtener todos los usuarios
	if err := r.db.WithContext(ctx).Find(&daoUsers).Error; err != nil {
		return nil, err
	}

	// Convertir de DAO a Domain
	domainUsers := make([]domain.UserDto, len(daoUsers))
	for i, daoUser := range daoUsers {
		// Desreferenciar si ToDomain() devuelve *domain.UserDto
		d := daoUser.ToDomain()
		if d == nil {
			domainUsers[i] = domain.UserDto{}
			continue
		}
		domainUsers[i] = *d
	}

	return domainUsers, nil
}

// Create inserta un nuevo usuario en MySQL
func (r *MySQLUsersRepository) Create(ctx context.Context, user domain.UserDto) (domain.UserDto, error) {
	daoUser := dao.FromDomain(user)

	// Usar GORM para crear el usuario
	if err := r.db.WithContext(ctx).Create(&daoUser).Error; err != nil {
		return domain.UserDto{}, err
	}

	// Retornar el usuario con el ID generado
	return daoUser.ToDomain(), nil
}

// GetByID busca un usuario por su ID
func (r *MySQLUsersRepository) GetByID(ctx context.Context, id int) (domain.UserDto, error) {
	var daoUser dao.UserModel

	// Usar GORM para buscar por ID
	if err := r.db.WithContext(ctx).First(&daoUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.UserDto{}, errors.New("user not found")
		}
		return domain.UserDto{}, err
	}

	return daoUser.ToDomain(), nil
}

// GetByEmail busca un usuario por su email
func (r *MySQLUsersRepository) GetByEmail(ctx context.Context, email string) (domain.UserDto, error) {
	var daoUser dao.UserModel

	// Usar GORM para buscar por email
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&daoUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.UserDto{}, errors.New("user not found")
		}
		return domain.UserDto{}, err
	}

	return daoUser.ToDomain(), nil
}

func (r *MySQLUsersRepository) GetUserByUsername(ctx context.Context, username string) (domain.UserDto, error) {
	var daoUser dao.UserModel
	// Usar GORM para buscar por username
	if err := r.db.Where("username = ?", username).First(&daoUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.UserDto{}, errors.New("user not found")
		}
		return domain.UserDto{}, err
	}

	return daoUser.ToDomain(), nil
}

// Update actualiza un usuario existente
func (r *MySQLUsersRepository) Update(ctx context.Context, id int, user domain.UserDto) (domain.UserDto, error) {
	daoUser := dao.FromDomain(user)
	daoUser.ID = id // Asegurar que el ID sea el correcto

	// Usar GORM para actualizar
	if err := r.db.WithContext(ctx).Save(&daoUser).Error; err != nil {
		return domain.UserDto{}, err
	}

	return daoUser.ToDomain(), nil
}

// Delete elimina un usuario por ID
func (r *MySQLUsersRepository) Delete(ctx context.Context, id int) error {
	// Usar GORM para eliminar
	result := r.db.WithContext(ctx).Delete(&dao.UserModel{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
