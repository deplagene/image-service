package broker

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
}

func NewKafka() *Kafka {
	return &Kafka{}
}

func (k *Kafka) Publish(ctx context.Context, topic string, data string) error {
	const op = "services.broker.kafka.Publish"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
	})

	defer func() {
		err := writer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(data),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (k *Kafka) Consume(ctx context.Context, topic string, groupId string) error {
	const op = "services.broker.kafka.Consume"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
		GroupID: groupId,
	})

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		fmt.Printf("message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
	}
	
	return nil
}
