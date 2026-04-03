package main

import (
	"mini-project-sanbercode/database"
	"mini-project-sanbercode/routers"
)

func main() {
	database.ConnectDB()

	routers.StartServer().Run(":8080")
}