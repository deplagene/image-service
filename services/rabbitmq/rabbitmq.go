package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMqClient(conn *amqp.Connection) (*Client, error) {
	ch, err := conn.Channel()
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		conn: conn,
		ch: ch,
	}, nil
}

func Connect(userName, password, host string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp::/%s:%s@%s", userName, password, host))
}

func(r *Client) Close() error {
	return r.ch.Close()
}

func(r *Client) CreateQueue(queueName string, durable, autoDelete bool) (amqp.Queue, error) {
	queue, err := r.ch.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}
	return queue, nil
}

func(r *Client) CreateExchange(name, typeOfExchange string, durable bool) error {
	return r.ch.ExchangeDeclare(name, typeOfExchange, durable, false, false, false, nil)
}

func(r *Client) CreateBinding(name, routingKey, exchangeName string) error {
	return r.ch.QueueBind(name, routingKey, exchangeName, false, nil)
}

func(r *Client)	Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	return r.ch.PublishWithContext(ctx, exchange, routingKey, false, false, options)
}

func(r *Client) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return r.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}