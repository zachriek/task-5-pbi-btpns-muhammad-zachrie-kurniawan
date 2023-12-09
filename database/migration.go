package database

import (
	"log"
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/models"
)

func Migrate() {
	Instance.AutoMigrate(&models.User{}, &models.Photo{})
	log.Println("Database Migration Success!")
}
