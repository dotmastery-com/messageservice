package main

import (
	"fmt"
	"net/http"
	"realtime-chat/services"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var appName = "User Service"

func main() {
	fmt.Printf("Starting %v\n", appName)

	startWebServer()
	startKafkaClient()

}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &services.Client{
		Id:     uuid.Must(uuid.NewV4()).String(),
		Socket: conn,
		Send:   make(chan []byte)}

	services.Manager.Register(client)

	go client.Read()
	go client.Write()
}

func startWebServer() {

	services.Manager = &services.WebSocketClientManager{
		BroadcastChannel:  make(chan []byte),
		RegisterChannel:   make(chan *services.Client),
		UnregisterChannel: make(chan *services.Client),
		Clients:           make(map[*services.Client]bool),
	}

	go services.Manager.Start()

	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)

}

func startKafkaClient() {

	services.Kafka = services.KafkaClient{
		Topic: "messages",
	}

	go services.Kafka.ConsumeTopic(true)
}
