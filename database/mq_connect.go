package database

import (
	"fmt"
	"tracker/app/config"

	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() (*amqp091.Connection, error) {
	env, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	buildConnectString := fmt.Sprintf("amqp://%s:%s@%s:%s/", env.MQ.Username, env.MQ.Password, env.MQ.Host, env.MQ.Port)
	fmt.Println(buildConnectString)
	conn, err := amqp091.Dial(buildConnectString)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	defer ch.Close()

	// Declare a fanout exchange for pub/sub
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	// Create a random, exclusive queue for this subscriber
	q, err := ch.QueueDeclare(
		"",    // empty name means random queue
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	// Bind queue to the exchange
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
