package routers

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hichoma.chat.dev/internal/handlers"
)

func InitializeRouter() {
	wsServer := socketio.NewServer(nil)
	wsServer.OnConnect("/", handlers.WSConnect)
	wsServer.OnDisconnect("/", handlers.WSDisconnect)

	go wsServer.Serve()
	defer wsServer.Close()
	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.POST("/register", handlers.Register)
	app.Any("/socket.io/*", func(context echo.Context) error {
		wsServer.ServeHTTP(context.Response(), context.Request())
		return nil
	})
	app.Logger.Fatal(app.Start(":5000"))
}
