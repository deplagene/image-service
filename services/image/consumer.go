package image

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"teach-stack/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	service types.ImageService
	broker types.MessageBroker
}

func NewConsumer(service types.ImageService, broker types.MessageBroker) *Consumer {
	return &Consumer{
		service: service,
		broker: broker,
	}
}

func (c *Consumer) Listen() error {
	if err := c.broker.CreateExchange("image", amqp.ExchangeFanout, true); err != nil {
		return fmt.Errorf("failed to create exchange: %w", err)
	}
	queue, err := c.broker.CreateQueue("", true, false)
    if err != nil {
        return fmt.Errorf("failed to create queue: %w", err)
    }

	if err := c.broker.CreateBinding(queue.Name, "", "image"); err != nil {
		return fmt.Errorf("failed to create binding: %w", err)
	}

	deliveries, err := c.broker.Consume(queue.Name, "", true)
    if err != nil {
        return fmt.Errorf("failed to consume queue: %w", err)
    }

	slog.Info("Consumer started, waiting for tasks...")

	forever := make(chan bool)
	go func() {
		for d := range deliveries {
			var task types.PhotoProcessingTask
            if err := json.Unmarshal(d.Body, &task); err != nil {
                slog.Error("Failed to unmarshal task", "error", err)
                continue
            }
			slog.Info("Received task", "fileName", task.FileName, "filters", task.Filters)
            if err := c.service.AddFilters(task); err != nil {
                slog.Error("Failed to process task", "error", err, "fileName", task.FileName)
                continue
            }
		}
	}()

	<- forever
	return nil
}