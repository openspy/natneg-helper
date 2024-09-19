package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"openspy.net/natneg-helper/src/Handlers"
	"openspy.net/natneg-helper/src/Messages"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	var core Handlers.NatNegCore

	var outboundHandler Handlers.IOutboundHandler
	var amqpHandler AMQPOutboundHandler = AMQPOutboundHandler{}
	outboundHandler = &amqpHandler

	portProbeDriver := os.Getenv("UNSOLICITED_PORT_PROBE_DRIVER")
	ipProbeDriver := os.Getenv("UNSOLICITED_IP_PROBE_DRIVER")
	ipPortProbeDriver := os.Getenv("UNSOLICITED_IPPORT_PROBE_DRIVER")

	core.Init(outboundHandler, 20, portProbeDriver, ipProbeDriver, ipPortProbeDriver)

	var amqpAddress string = os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(amqpAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"natneg-core", // name
		false,         // durable
		true,          // delete when unused
		true,          // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	amqpHandler.amqpQueue = q
	amqpHandler.amqpCtx = ctx
	amqpHandler.amqpChannel = ch

	ch.QueueBind(q.Name, "natneg.core", "openspy.natneg", false, nil)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for {
			core.Tick()
		}
	}()

	go func() {
		for d := range msgs {
			var msg Messages.Message
			err = json.Unmarshal(d.Body, &msg)
			if err == nil {
				Handlers.HandleMessage(core, outboundHandler, msg)
			}

		}
	}()

	<-forever
}
