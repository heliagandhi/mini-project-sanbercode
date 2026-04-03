package main

import (
	"mini-project-sanbercode/database"
	"mini-project-sanbercode/routers"
	"os"
)

func main() {
	database.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	routers.StartServer().Run(":" + port)
}