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

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.POST("/register", handlers.Register)
	app.POST("/login", handlers.Login)
	app.GET("/checkToken", handlers.CheckToken)

	wsProtectedRoutes := app.Group("")
	wsProtectedRoutes.Use(middlewareFc.WsAuthentificate)
	wsProtectedRoutes.GET("/ws", server.WSEndpoint)

	protectedRoutes := app.Group("")
	protectedRoutes.Use(middlewareFc.Authentificate)
	protectedRoutes.GET("/messages/:friend", handlers.GetUserChat)
	protectedRoutes.POST("/friends", handlers.AddUserFriend)
	protectedRoutes.GET("/friends", handlers.GetUserFriendList)

	app.Logger.Fatal(app.Start(":5000"))
}
