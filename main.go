package main

import (
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/database"
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/router"
)

func main() {
	database.Connect()
	database.Migrate()

	router := router.Setup()

	router.Run(":8000")
}
