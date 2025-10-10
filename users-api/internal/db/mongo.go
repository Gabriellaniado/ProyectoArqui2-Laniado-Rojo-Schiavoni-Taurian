package db

import (
	"users-api/internal/config"
	"users-api/internal/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	// Configurar GORM con logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Conectar a MySQL usando el DSN
	db, err := gorm.Open(mysql.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return nil, err
	}

	// Auto-migrar los modelos (crear tablas si no existen)
	err = db.AutoMigrate(&dao.UserModel{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
