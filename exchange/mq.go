package exchange

import (
	"github.com/streadway/amqp"
	"encoding/json"
)

const DefaultTTL = Min * 1000

type MQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func (mq *MQ) Init() {

	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		LogErr.Error("Failed to connect to RabbitMQ:" + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		LogErr.Error("Failed to open a channel:" + err.Error())
	}

	args := amqp.Table{"x-message-ttl": int32(DefaultTTL)}
	q, err := ch.QueueDeclare(
		"kline", // name
		false,   // durable
		true,    // delete when unused
		false,   // exclusive
		false,   // no-wait
		args,    // arguments
	)
	if err != nil {
		LogErr.Error("Failed to declare a queue:" + err.Error())
	}

	mq.conn = conn
	mq.channel = ch
	mq.queue = &q
}

func (mq *MQ) send(msg interface{}) {

	b, err := json.Marshal(msg)
	if err != nil {
		LogErr.Error("MQ marshal error:" + err.Error())
		return
	}

	err = mq.channel.Publish(
		"",
		mq.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		})

	if err != nil {
		LogErr.Error("MQ send error:" + err.Error())
	}
}

func (mq *MQ) close() {
	mq.conn.Close()
	mq.conn.Close()
}
