package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/alexohneander/GoZilla/pkg/model"
)

func Migrate() {
	// Define DB
	db, err := GetDB()
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Peer{})
}

func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gozilla.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
