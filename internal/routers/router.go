package routers

import (
	"github.com/labstack/echo/v4"
	"hichoma.chat.dev/internal/handlers"
)

func InitializeRouter() {
	app := echo.New()

	app.POST("/register", handlers.Register)
	app.Logger.Fatal(app.Start(":5000"))
}
