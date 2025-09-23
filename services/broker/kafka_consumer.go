package broker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

func NewConsumer(brokers []string, topic, groupId string) (*Consumer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupId,
		CommitInterval: 0,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
	})

	return &Consumer{
		reader: reader,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (c *Consumer) Consume(handler func(ctx context.Context, msg kafka.Message) error) {
	const op = "broker.kafka_consumer.Consumer.Consume"

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				log.Printf("%s: consumer is shutting down...", op)
				return
			default:
				msg, err := c.reader.FetchMessage(c.ctx)
				if err != nil {
					if c.ctx.Err() != nil {
						return
					}
					log.Printf("%s: %v\n", op, err)
					continue
				}
				if err := handler(c.ctx, msg); err != nil {
					log.Printf("%s: error processing message (offset %d): %v\n", op, msg.Offset, err)
				} else {
					if err := c.reader.CommitMessages(c.ctx, msg); err != nil {
						log.Printf("%s: error committing message (offset %d): %v\n", op, msg.Offset, err)
					}
				}
			}
		}
	}()

}

func (c *Consumer) Close() error {
	const op = "broker.kafka_consumer.Consumer.Close"

	c.cancel()
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
