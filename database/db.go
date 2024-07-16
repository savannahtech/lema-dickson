package database

import (
	"log"

	"github.com/midedickson/github-service/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectToDB() {
	d, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database sucessfully")
	DB = d
}

func AutoMigrate() {
	log.Println("Auto Migrating Models...")
	err := DB.AutoMigrate(&models.Repository{}, &models.Commit{})
	if err != nil {
		panic(err)
	}
	log.Println("Migrated DB Successfully")
}
