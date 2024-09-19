package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"net/netip"

	amqp "github.com/rabbitmq/amqp091-go"
	"openspy.net/natneg-helper/src/Handlers"
	"openspy.net/natneg-helper/src/Messages"
)

type AMQPOutboundHandler struct {
	amqpCtx     context.Context
	amqpQueue   amqp.Queue
	amqpChannel *amqp.Channel
}

func (h *AMQPOutboundHandler) SendMessage(msg Messages.Message) {
	body, jsonErr := json.Marshal(msg)

	if jsonErr != nil {
		log.Fatalf("Failed to marshal msg\n")
		return
	}

	err := h.amqpChannel.PublishWithContext(h.amqpCtx,
		"openspy.natneg",  // exchange
		"natneg.endpoint", // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to publish message: %s\n", err.Error())
		return
	}
}

func (h *AMQPOutboundHandler) SendDeadbeatMessage(client *Handlers.NatNegSessionClient) {
	var msg Messages.Message
	var connectMsg Messages.ConnectMessage

	connectMsg.RemoteIP = 0
	connectMsg.RemotePort = 0

	connectMsg.Finished = Messages.FINISHED_ERROR_DEAD_PARTNER
	if client.GotRemoteData() {
		connectMsg.GotYourData = 1
	} else {
		connectMsg.GotYourData = 0
	}

	msg.Type = "connect"
	msg.Cookie = client.Cookie
	msg.Version = client.Version

	for _, info := range client.InitAddresses {
		msg.DriverAddress = info.DriverAddress.String()
		msg.Address = info.Address.String()
		msg.Message = connectMsg
		h.SendMessage(msg)

	}
}

func (h *AMQPOutboundHandler) SendConnectMessage(client *Handlers.NatNegSessionClient, ipAddress netip.AddrPort) {
	var msg Messages.Message
	var connectMsg Messages.ConnectMessage

	var ipBuff = ipAddress.Addr().As4()

	cpyBuff := make([]byte, 4)
	cpyBuff[3] = ipBuff[0]
	cpyBuff[2] = ipBuff[1]
	cpyBuff[1] = ipBuff[2]
	cpyBuff[0] = ipBuff[3]

	ipAddr := binary.BigEndian.Uint32(cpyBuff)
	connectMsg.RemoteIP = int(ipAddr)
	connectMsg.RemotePort = int(ipAddress.Port())

	connectMsg.Finished = Messages.FINISHED_NOERROR
	if client.GotRemoteData() {
		connectMsg.GotYourData = 1
	} else {
		connectMsg.GotYourData = 0
	}

	msg.Type = "connect"
	msg.Cookie = client.Cookie
	msg.Version = client.Version

	for _, info := range client.InitAddresses {
		msg.DriverAddress = info.DriverAddress.String()
		msg.Address = info.Address.String()
		msg.Message = connectMsg
		h.SendMessage(msg)

	}
}
