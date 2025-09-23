package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        true,
		WriteTimeout: 10 * time.Second,
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Publish(ctx context.Context, msg kafka.Message) error {
	const op = "broker.kafka_producer.Producer.Publish"

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Producer) Close() error {
	const op = "broker.kafka_producer.Producer.Close"

	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
