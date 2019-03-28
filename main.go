package main

import (
	"fmt"
	"net/http"
	"realtime-chat/services"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func main() {
	fmt.Println("Starting application...")

	services.Manager = services.ClientManager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *services.Client),
		Unregister: make(chan *services.Client),
		Clients:    make(map[*services.Client]bool),
	}

	services.Kafka = services.KafkaClient{
		Topic: "test",
	}

	//go services.Kafka.Connect()
	go services.Kafka.ConnectToTopic()
	go services.Manager.Start()

	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &services.Client{Id: uuid.Must(uuid.NewV4(), error).String(), Socket: conn, Send: make(chan []byte)}

	services.Manager.Register <- client

	go client.Read()
	go client.Write()
}
