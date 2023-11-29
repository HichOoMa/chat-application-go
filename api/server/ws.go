package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Msg  chan []byte
}

type WSServer struct {
	Conns map[string]*Client
	mu    sync.Mutex
}

func StartWS() *WSServer {
	return &WSServer{
		Conns: make(map[string]*Client),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (server *WSServer) WSEndpoint(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	client := server.joinClient(ws)
	defer server.disconnect(client)

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			ctx.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			ctx.Logger().Error(err)
			return err
		}
		response := fmt.Sprintf("Hello %s\n", msg)
		ws.WriteMessage(websocket.TextMessage, []byte(response))
	}
}

func (server *WSServer) joinClient(ws *websocket.Conn) *Client {
	client := Client{
		ID:   uuid.Must(uuid.NewRandom()).String(),
		Conn: ws,
		Msg:  make(chan []byte),
	}
	server.mu.Lock()
	defer server.mu.Unlock()
	// add to the connected user list
	server.Conns[client.ID] = &client

	return &client
}

func (server *WSServer) disconnect(client *Client) {
	defer client.Conn.Close()
	delete(server.Conns, client.ID)
}
