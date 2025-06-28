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

	var skipERT bool = false

	_, hasKey := os.LookupEnv("SKIP_ERT")
	if hasKey {
		skipERT = true
	}

	core.Init(outboundHandler, 10, portProbeDriver, ipProbeDriver, ipPortProbeDriver, skipERT)

	//attempt outbound amqp connection max 20 tries
    var attempts = 20;
    var outboundConn amqp.Connection
    var err error
    var amqpAddress string = os.Getenv("RABBITMQ_URL")
    for attempts > 0 {
        outboundConn, err := amqp.Dial(amqpAddress)
        if err != nil {
           attempts -= 1
       } else {
           failOnError(err, "Failed to connect to RabbitMQ")
           defer outboundConn.Close()
       }
       // Sleep for attempts * 5 secs
       var duration = (20 - attempts) * 5
       time.Sleep(time.Duration(duration)  * time.Second)
    }

	outboundChannel, err := outboundConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer outboundChannel.Close()

	outbountCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	amqpHandler.amqpCtx = outbountCtx
	amqpHandler.amqpChannel = outboundChannel

	//make listener connection, etc
	listenConn, err := amqp.Dial(amqpAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer listenConn.Close()

	chListen, err := listenConn.Channel()
	failOnError(err, "Failed to open a channel")
	defer chListen.Close()

	q, err := chListen.QueueDeclare(
		"natneg-core", // name
		false,         // durable
		true,          // delete when unused
		true,          // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	chListen.QueueBind(q.Name, "natneg.core", "openspy.natneg", false, nil)

	msgs, err := chListen.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for {
			core.Tick()
		}
	}()

	for d := range msgs {
		var msg Messages.Message
		err = json.Unmarshal(d.Body, &msg)
		if err == nil {
			Handlers.HandleMessage(core, outboundHandler, msg)
		} else {
			log.Printf("Failed to unmarshal msg\n")
		}
	}
}
