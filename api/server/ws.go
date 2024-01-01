package server

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"hichoma.chat.dev/api/handlers"
	"hichoma.chat.dev/internal/models"
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
	client := server.joinClient(ws, ctx)
	defer server.disconnect(client)

	for {
		msg := models.WsMessage{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			ctx.Logger().Error(err.Error())
			break
		}

		err = server.sendMessage(&msg, ctx)
		if err != nil {
			ctx.Logger().Error("message sent failed")
		}
	}
	return nil
}

func (server *WSServer) joinClient(ws *websocket.Conn, ctx echo.Context) *Client {
	userId := ctx.Request().Header.Values("user")
	client := Client{
		ID:   strings.Join(userId, ""),
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

func (server *WSServer) sendMessage(msg *models.WsMessage, ctx echo.Context) error {
	msgResponse := handlers.AddNewMessage(msg, ctx)
	server.Conns[msgResponse.SenderID].Conn.WriteJSON(msgResponse)

	if receiverClient, ok := server.Conns[msgResponse.ReceiverID]; ok {

		err := receiverClient.Conn.WriteJSON(msgResponse)
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}
