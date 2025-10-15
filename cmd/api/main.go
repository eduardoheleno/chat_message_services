package main

import (
	"chat_service/cmd/api/database"
	"chat_service/cmd/api/route"
)

func main() {
	db := database.InitDatabase()

	router := route.InitRoutes(db)
	router.Run()
}
