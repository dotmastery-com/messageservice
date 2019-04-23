package services

import (
	"encoding/json"
	"fmt"
	"log"
	"realtime-chat/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/xid"
)

type KafkaClient struct {
	Topic string
}

var LocationID = xid.New().String()

func (kc *KafkaClient) ConsumeTopic(ignoreMessagesFromSource bool, numMessages int) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"client.id":         "client",
		"auto.offset.reset": "latest",
	})

	if err != nil {
		panic(err)
	}

	log.Println("Consuming from Topic", kc.Topic)

	c.Subscribe(kc.Topic, nil)
	c.Poll(0)

	n := 0
	for n < numMessages {

		message := model.Message{}

		msg, err := c.ReadMessage(-1)
		n++

		if err != nil {

			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)

		} else {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			err = json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Printf("Unable to parse messages %s = %s\n", string(msg.Key), string(msg.Value))
			}

			//-- we may choose to ignore messages that this location has posted to Kafka since they may have already been broadcast
			if ignoreMessagesFromSource {
				if message.SourceLocation != LocationID {
					Manager.Broadcast(msg.Value)

				} else {
					log.Printf("Ignoring message from source location. \n")

				}
			} else {
				Manager.Broadcast(msg.Value)
			}

		}
	}

	c.Close()

}

func (kc *KafkaClient) SendMessage(message model.Message) error {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	message.Id = xid.New().String()
	message.SourceLocation = LocationID
	jsonBytes, _ := json.Marshal(message)

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kc.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(jsonBytes),
		Key:            []byte(message.Id),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(1 * 1000)

	return nil
}
