package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hichoma.chat.dev/api/handlers"
	middlewareFc "hichoma.chat.dev/api/middleware"
)

func InitializeServer() {
	app := echo.New()

	server := StartWS()

	app.Use(middleware.Logger())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.POST("/register", handlers.Register)
	app.POST("/login", handlers.Login)
	app.GET("/checkToken", handlers.CheckToken)

	protectedRoutes := app.Group("/")
	protectedRoutes.Use(middlewareFc.Authentificate)
	protectedRoutes.GET("/ws", server.WSEndpoint)

	app.Logger.Fatal(app.Start(":5000"))
}
