package types

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Image struct {
	Id          uuid.UUID
	PayloadName string
	Payload     io.Reader
	Size        int64
	Url         string
}

type ImageCreatedEvent struct {
	Url string
}

type S3_Storage interface {
	Connect() error
	Upload(ctx context.Context, img Image) (string, error)
	GetTemplUrl(ctx context.Context, fileName string) (string, error)
	DownloadImage(ctx context.Context, bucketName string, fileName string) (*minio.Object, error)
}

type MessageBroker interface {
	Close() error
	CreateQueue(queueName string, durable, autoDelete bool) (amqp.Queue, error)
	CreateExchange(name, typeOfExchange string, durable bool) error
	CreateBinding(name, routingKey, exchangeName string) error
	Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error
	Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
}

type ImageStore interface {
	Create(Image) error
	GetByUrl(url string) (*Image, error)
	GetById(id uuid.UUID) (*Image, error)
	Update(id uuid.UUID, url string) error
}

type ImageService interface {
	AddFilters(task PhotoProcessingTask) error
}

type PhotoProcessingTask struct {
	Id       uuid.UUID `json:"id"`
	FileName string    `json:"file_name"`
	Filters  []string  `json:"filters"`
}
