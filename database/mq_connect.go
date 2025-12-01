package database

import (
	"fmt"
	"tracker/app/config"

	"github.com/rabbitmq/amqp091-go"
)

type MQChannels struct {
	Connection *amqp091.Connection
	Channels   map[string]*amqp091.Channel // Keyed by task type
}

func setupTaskQueue(conn *amqp091.Connection, routingKey string, queueName string) (*amqp091.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	exchange := "ai_events_exchange"
	err = ch.ExchangeDeclare(
		exchange, "direct", true, false, false, false, nil,
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	err = ch.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return ch, nil
}

func ConnectRabbitMQ() (*MQChannels, error) {
	env, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	url := env.MQ.AmazonMQ
	if url == "" {
		url = fmt.Sprintf("amqp://%s:%s@%s:%s/", env.MQ.Username, env.MQ.Password, env.MQ.Host, env.MQ.Port)
	} else {
		url = fmt.Sprintf("amqps://%s:%s@%s", env.MQ.Username, env.MQ.Password, env.MQ.AmazonMQ)
	}
	fmt.Println(url)
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	fmt.Println("RabbitMQ connected.")
	channels := make(map[string]*amqp091.Channel)
	// Sentiment analysis
	ch, err := setupTaskQueue(conn, "generate", "worker")
	if err != nil {
		conn.Close()
		return nil, err
	}
	channels["generate"] = ch
	// student reports
	chEmailSend, err := setupTaskQueue(conn, "report", "micro_report")
	if err != nil {
		conn.Close()
		return nil, err
	}
	channels["report"] = chEmailSend
	// Org reports
	pgReportSender, err := setupTaskQueue(conn, "pgdata", "micro_pgreport")
	if err != nil {
		conn.Close()
		return nil, err
	}
	channels["pgdata"] = pgReportSender
	// AI Quiz and materials generation
	chMaterials, err := setupTaskQueue(conn, "produce", "micro_materials")
	if err != nil {
		conn.Close()
		return nil, err
	}
	channels["produce"] = chMaterials

	cGrader, err := setupTaskQueue(conn, "grader", "micro_grader")
	if err != nil {
		conn.Close()
		return nil, err
	}
	channels["grader"] = cGrader

	return &MQChannels{
		Connection: conn,
		Channels:   channels,
	}, nil
}
