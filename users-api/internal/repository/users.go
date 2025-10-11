package repository

import (
	"context"
	"errors"
	"strings"
	"users-api/internal/dao"
	"users-api/internal/domain"

	"gorm.io/gorm"
)

// UsersRepository define las operaciones de datos para Users
type UsersRepository interface {
	List(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByID(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, id int, user domain.User) (domain.User, error)
	Delete(ctx context.Context, id int) error
}

// MySQLUsersRepository implementa UsersRepository usando MySQL con GORM
type MySQLUsersRepository struct {
	db *gorm.DB
}

func NewMySQLUsersRepository(db *gorm.DB) UsersRepository {
	return &MySQLUsersRepository{db: db}
}

// List obtiene todos los usuarios de MySQL
func (r *MySQLUsersRepository) List(ctx context.Context) ([]domain.User, error) {
	var daoUsers []dao.UserModel

	if err := r.db.WithContext(ctx).Find(&daoUsers).Error; err != nil {
		return nil, err
	}

	// Convertir de DAO a Domain - corregido sin desreferenciación
	domainUsers := make([]domain.User, len(daoUsers))
	for i, daoUser := range daoUsers {
		domainUsers[i] = daoUser.ToDomain()
	}

	return domainUsers, nil
}

// Create inserta un nuevo usuario en MySQL
func (r *MySQLUsersRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	daoUser := dao.FromDomain(user)

	if err := r.db.WithContext(ctx).Create(&daoUser).Error; err != nil {
		return domain.User{}, err
	}

	return daoUser.ToDomain(), nil
}

// GetByID busca un usuario por su ID
func (r *MySQLUsersRepository) GetByID(ctx context.Context, id int) (domain.User, error) {
	var daoUser dao.UserModel

	if err := r.db.WithContext(ctx).First(&daoUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return daoUser.ToDomain(), nil
}

// GetByEmail busca un usuario por su email
func (r *MySQLUsersRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var daoUser dao.UserModel

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&daoUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return daoUser.ToDomain(), nil
}

// Update actualiza un usuario existente
/*func (r *MySQLUsersRepository) Update(ctx context.Context, id int, user domain.User) (domain.User, error) {
	daoUser := dao.FromDomain(user)
	daoUser.ID = id

	if err := r.db.WithContext(ctx).Omit("created_at").Save(&daoUser).Error; err != nil {
		return domain.User{}, err
	}

	return daoUser.ToDomain(), nil
}*/
// Update actualiza un usuario existente
func (r *MySQLUsersRepository) Update(ctx context.Context, id int, user domain.User) (domain.User, error) {
	// Preparar mapa de campos a actualizar
	updates := map[string]interface{}{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"is_admin":   user.IsAdmin,
	}

	// Solo agregar password si no está vacío
	if strings.TrimSpace(user.Password) != "" {
		updates["password_hash"] = user.Password
	}

	// Actualizar solo los campos especificados
	result := r.db.WithContext(ctx).
		Model(&dao.UserModel{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	// Obtener el usuario actualizado
	var userDAO dao.UserModel
	r.db.WithContext(ctx).Where("id = ?", id).First(&userDAO)

	return userDAO.ToDomain(), nil
}

// Delete elimina un usuario por ID
func (r *MySQLUsersRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&dao.UserModel{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
