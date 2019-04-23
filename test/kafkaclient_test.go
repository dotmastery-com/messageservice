package test

import (
	"encoding/json"
	"fmt"
	"realtime-chat/model"
	"realtime-chat/services"
	"strconv"
	"testing"
	"time"
)

type MockClientManager struct {
	messages []model.Message
}

func (m *MockClientManager) Broadcast(message []byte) {

	model := model.Message{}

	json.Unmarshal(message, &model)

	m.messages = append(m.messages, model)
	fmt.Printf("Broadcast called %v", model)
}

func (m *MockClientManager) Start() {
	fmt.Println("Start called")

}

func (m *MockClientManager) Unregister(c *services.Client) {
	fmt.Println("Unregister called")
}

func (m *MockClientManager) Register(c *services.Client) {
	fmt.Println("Register called")
}

/*
  Test that when we send (produce) a message to Kafka we receive (consume) the message with the correct content
*/
func TestConsumingMessageCorrectContentReturned(t *testing.T) {

	mockClientManager := &MockClientManager{
		messages: make([]model.Message, 0),
	}

	services.Manager = mockClientManager

	services.Kafka = services.KafkaClient{
		Topic: "test",
	}

	message := model.Message{
		Content: strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	services.Kafka.SendMessage(message)

	//-- test that we get the correct content
	services.Kafka.ConsumeTopic(false, 1)

	if mockClientManager.messages[0].Content != message.Content {
		t.Errorf("Unexpected message! Got %v want %v", message.Content, mockClientManager.messages[0].Content)
	}

}

/*
  Tests that messages sent from the local location are not re-broadcast when they're consumed from Kafka.
  Messages that come through the local instance are automatically broadcast to attached web sockets before being published to Kafka
  If we pick up messages that we've already sent, then the users will see duplicates
*/
func TestConsumingMessageSourceMessagesIgnored(t *testing.T) {

	mockClientManager := &MockClientManager{}

	services.Manager = mockClientManager

	services.Kafka = services.KafkaClient{
		Topic: "test",
	}

	message := model.Message{
		Content: strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	services.Kafka.SendMessage(message)

	//-- test that we ignore the message as it comes from the local source
	services.Kafka.ConsumeTopic(true, 1)

	if mockClientManager.messages != nil {
		t.Errorf("Unexpected message! Got %v want %v", message.Content, "")
	}

}

/*
  Test that we can handle multiple messages, and that they're broadcast in the correct order
*/
func TestConsumingMultipleMessages(t *testing.T) {

	mockClientManager := &MockClientManager{}

	services.Manager = mockClientManager

	services.Kafka = services.KafkaClient{
		Topic: "test",
	}

	var messages []model.Message

	//-- write 10 messages, each with a unique timestamp
	for i := 0; i < 10; i++ {
		message := model.Message{
			Content: strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		messages = append(messages, message)

		services.Kafka.SendMessage(message)

	}

	//-- consume 10 messages
	services.Kafka.ConsumeTopic(false, 10)

	//-- verify that the messages are received in the order they're sent
	for i := range messages {
		if messages[i].Content != mockClientManager.messages[i].Content {
			t.Errorf("Unexpected message! Got %v want %v", messages[i].Content, mockClientManager.messages[i].Content)
		}
	}

}

/*
 Test sending a message, and that it gets set with the correct location ID
*/
func TestSendingMessage(t *testing.T) {
	mockClientManager := &MockClientManager{}

	services.Manager = mockClientManager

	services.Kafka = services.KafkaClient{
		Topic: "test",
	}

	message := model.Message{
		Content: "Hello world",
	}

	services.Kafka.SendMessage(message)

	//-- consume 1 message
	services.Kafka.ConsumeTopic(false, 1)

	//-- verify that the location ID is what we expect
	if services.LocationID != mockClientManager.messages[0].SourceLocation {
		t.Errorf("Unexpected location ID! Got %v want %v", mockClientManager.messages[0].SourceLocation, services.LocationID)
	}

	//-- verify that the payload is correct
	if message.Content != mockClientManager.messages[0].Content {
		t.Errorf("Unexpected message! Got %v want %v", mockClientManager.messages[0].Content, message.Content)
	}

}
