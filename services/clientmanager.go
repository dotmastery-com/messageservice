package services

import (
	"encoding/json"
	"fmt"
	"realtime-chat/model"

	"github.com/gorilla/websocket"
)

type ClientManager interface {
	Broadcast(message []byte)
	Start()
	Unregister(c *Client)
	Register(c *Client)
}

var Manager ClientManager

var Kafka KafkaClient

type WebSocketClientManager struct {
	Clients           map[*Client]bool
	BroadcastChannel  chan []byte
	RegisterChannel   chan *Client
	UnregisterChannel chan *Client
}

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
}

func (manager *WebSocketClientManager) Broadcast(message []byte) {

	manager.send(message, nil)
}

func (manager *WebSocketClientManager) Unregister(c *Client) {

	manager.UnregisterChannel <- c
}

func (manager *WebSocketClientManager) Register(c *Client) {

	manager.RegisterChannel <- c
}

func (manager *WebSocketClientManager) Start() {
	for {
		select {
		case conn := <-manager.RegisterChannel:
			manager.Clients[conn] = true
			jsonMessage, _ := json.Marshal(&model.Message{Content: "/A new user has connected."})
			manager.send(jsonMessage, conn)
		case conn := <-manager.UnregisterChannel:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&model.Message{Content: "/A user has disconnected."})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.BroadcastChannel:
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

func (manager *WebSocketClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister(c)

		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister(c)
			c.Socket.Close()
			break
		}

		messagea := model.Message{Content: string(message)}
		jsonMessage, _ := json.Marshal(&messagea)

		Manager.Broadcast(jsonMessage)
		err = Kafka.SendMessage(messagea)

		if err != nil {
			fmt.Println("Error writing to Kafka.")
		}

	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
