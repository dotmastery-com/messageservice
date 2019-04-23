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

func TestConsumingMultipleMessages(t *testing.T) {

	mockClientManager := &MockClientManager{}

	services.Manager = mockClientManager

	services.Kafka = services.KafkaClient{
		Topic: "test",
	}

	var messages []model.Message

	//-- write 10 messages
	for i := 0; i < 10; i++ {
		message := model.Message{
			Content: strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		messages = append(messages, message)

		services.Kafka.SendMessage(message)

	}

	//-- consume 10 messages
	services.Kafka.ConsumeTopic(false, 10)

	for i := range messages {
		if messages[i].Content != mockClientManager.messages[i].Content {
			t.Errorf("Unexpected message! Got %v want %v", messages[i].Content, mockClientManager.messages[i].Content)
		}
	}

}

func TestSendingMessage(t *testing.T) {
	services.Kafka = services.KafkaClient{
		Topic: "test",
	}

	message := model.Message{

		Content:        "Worald",
		SourceLocation: "Here",
	}

	services.Kafka.SendMessage(message)

}
