package main

import (
	"hichoma.chat.dev/api/server"
	"hichoma.chat.dev/internal/config"
	"hichoma.chat.dev/internal/database"
)

func main() {
	config.InitializeAppConfig()
	database.InitializeDatabase()
	server.InitializeServer()
}
