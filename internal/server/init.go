package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hichoma.chat.dev/internal/handlers"
	middlewareFc "hichoma.chat.dev/internal/middleware"
)

func InitializeServer() {
	app := echo.New()

	// app.Use(middleware.Logger())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.POST("/register", handlers.Register)
	protectedRoutes := app.Group("/")
	protectedRoutes.Use(middlewareFc.Authentificate)
	protectedRoutes.GET("/ws", handlers.WSEndpoint)
	app.Logger.Fatal(app.Start(":5000"))
}
