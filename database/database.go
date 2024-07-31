package database

import (
	"github.com/iadityanath8/gopi/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDriver() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return err
	}

	DB = db
	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.User{}, &models.CartItem{}, &models.Cart{})
	if err != nil {
		return err
	}
	return nil
}
