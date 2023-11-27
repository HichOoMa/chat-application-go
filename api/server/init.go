package server

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hichoma.chat.dev/api/handlers"
	middlewareFc "hichoma.chat.dev/api/middleware"
)

var WSConnections = make(map[string]websocket.Conn)

func InitializeServer() {
	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.POST("/register", handlers.Register)
	app.POST("/login", handlers.Login)
	app.GET("/checkToken", handlers.CheckToken)

	protectedRoutes := app.Group("/")
	protectedRoutes.Use(middlewareFc.Authentificate)
	protectedRoutes.GET("/ws", handlers.WSEndpoint)

	app.Logger.Fatal(app.Start(":5000"))
}
