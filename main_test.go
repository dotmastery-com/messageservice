package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"realtime-chat/model"
	"realtime-chat/services"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TestSimpleSendReceiveSingleClient(t *testing.T) {

	services.Manager = &services.WebSocketClientManager{
		BroadcastChannel:  make(chan []byte),
		RegisterChannel:   make(chan *services.Client),
		UnregisterChannel: make(chan *services.Client),
		Clients:           make(map[*services.Client]bool),
	}

	go services.Manager.Start()

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(WSPage))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	for i := 0; i < 10; i++ {
		if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
			t.Fatalf("%v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		var message model.Message
		err = json.Unmarshal(p, &message)
		fmt.Printf("%+v", message)

		expectedMessage := model.Message{
			Id:             "",
			Sender:         "",
			Recipient:      "",
			Content:        "hello",
			SourceLocation: "",
		}

		if message != expectedMessage {
			t.Fatalf("bad message")
		}

	}

}

func TestSimpleSendReceiveMultipleClient(t *testing.T) {

	services.Manager = &services.WebSocketClientManager{
		BroadcastChannel:  make(chan []byte),
		RegisterChannel:   make(chan *services.Client),
		UnregisterChannel: make(chan *services.Client),
		Clients:           make(map[*services.Client]bool),
	}

	go services.Manager.Start()

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(WSPage))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect Client one to the server
	wsOne, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer wsOne.Close()

	// Connect client two to the server
	wsTwo, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer wsTwo.Close()

	//-- check that client one is notified of the new user
	_, p, err := wsOne.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	expectedMessage := model.Message{
		Id:             "",
		Sender:         "",
		Recipient:      "",
		Content:        "/A new user has connected.",
		SourceLocation: "",
	}

	if !compare(p, expectedMessage) {
		t.Fatalf("Incorrect message. Expected %v, got %v", expectedMessage, p)
	}

	// Send messages to the server from client one, test that both client one and client two get the messages
	for i := 0; i < 10; i++ {
		if err := wsOne.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
			t.Fatalf("%v", err)
		}

		expectedMessage := model.Message{
			Id:             "",
			Sender:         "",
			Recipient:      "",
			Content:        "hello",
			SourceLocation: "",
		}

		_, clientTwoMsg, err := wsTwo.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if !compare(clientTwoMsg, expectedMessage) {
			t.Fatalf("Incorrect message. Expected %v, got %v", expectedMessage, clientTwoMsg)
		}

		_, clientOneMsg, err := wsOne.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if !compare(clientOneMsg, expectedMessage) {
			t.Fatalf("Incorrect message. Expected %v, got %v", expectedMessage, clientOneMsg)
		}

	}

	//-- disconnect client two
	wsTwo.Close()

	//-- check that client one is notified of user leaving
	_, p, err = wsOne.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	expectedMessage = model.Message{
		Id:             "",
		Sender:         "",
		Recipient:      "",
		Content:        "/A user has disconnected.",
		SourceLocation: "",
	}

	if !compare(p, expectedMessage) {
		t.Fatalf("Incorrect message. Expected %v, got %v", expectedMessage, p)
	}

}

func compare(incomingMessage []byte, expectedMessage model.Message) bool {

	var message model.Message
	_ = json.Unmarshal(incomingMessage, &message)
	fmt.Printf("%+v", message)
	return message == expectedMessage

}
