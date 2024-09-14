package main

import (
	"log"

	"openspy.net/natneg-helper/src/Messages"
)

type AMQPOutboundHandler struct {
}

func (h *AMQPOutboundHandler) SendMessage(msg Messages.Message) {
	log.Printf("send msg\n")
}
