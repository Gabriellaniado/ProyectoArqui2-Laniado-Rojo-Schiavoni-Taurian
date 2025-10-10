package dao

import (
	"time"
)

type UserModel struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`          //PK
	Email        string    `gorm:"unique;not null;type:varchar(100)"` //Unique email
	PasswordHash string    `gorm:"longtext"`                          //Password Hash
	FirstName    string    `gorm:"type:varchar(100);not null"`
	LastName     string    `gorm:"type:varchar(100);not null"`
	IsAdmin      bool      `gorm:"default:false"` //Admin
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// ToDomain convierte de modelo MySQL a modelo de negocio
func (u UserModel) ToDomain() domain.UserDto {
	return domain.UserDto{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromDomain convierte de modelo de negocio a modelo MySQL
func FromDomain(domainUser domain.UserDto) UserModel {
	return UserModel{
		ID:           domainUser.ID,
		Email:        domainUser.Email,
		PasswordHash: domainUser.Password, // En domain se guarda como Password, en DAO como PasswordHash
		FirstName:    domainUser.FirstName,
		LastName:     domainUser.LastName,
		IsAdmin:      domainUser.IsAdmin,
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
	}
}
