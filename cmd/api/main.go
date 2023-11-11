package main

import (
	"hichoma.chat.dev/internal/config"
	"hichoma.chat.dev/internal/database"
	"hichoma.chat.dev/internal/routers"
)

func main() {
	config.InitializeAppConfig()
	database.InitializeDatabase()
	routers.InitializeRouter()
}
