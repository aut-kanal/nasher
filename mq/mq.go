package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/nasher/configuration"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel

	qAccepts amqp.Queue

	msgs <-chan amqp.Delivery
)

func SubscribeAcceptedMsgs(callback func(amqp.Delivery)) {
	go func() {
		for msg := range msgs {
			go callback(msg)
		}
	}()
}

func InitMessageQueue() {
	// Connection
	var err error
	conn, err = amqp.Dial(configuration.GetInstance().GetString("rabbit-mq.url"))
	if err != nil {
		logrus.WithError(err).Fatalln("can't connect to message queue")
	}

	// Channel
	channel, err = conn.Channel()
	if err != nil {
		logrus.WithError(err).Fatalln("can't create mq channel")
	}

	err = channel.ExchangeDeclare(
		configuration.GetInstance().GetString("rabbit-mq.accept-ex-name"), // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't declare accept exchange")
	}

	err = channel.ExchangeDeclare(
		configuration.GetInstance().GetString("rabbit-mq.reject-ex-name"), // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't declare reject exchange")
	}

	// Queue
	qAccepts, err = channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatalln("can't create accepts queue")
	}

	err = channel.QueueBind(
		qAccepts.Name, // queue name
		"",            // routing key
		configuration.GetInstance().GetString("rabbit-mq.accept-ex-name"), // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't bind queue to exchange")
	}

	// Consumer
	msgs, err = channel.Consume(
		qAccepts.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't init msg consumer")
	}

	logrus.Info("message queue initialized")
}

func Close() {
	conn.Close()
}
