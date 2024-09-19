package Handlers

import (
	"log"

	"openspy.net/natneg-helper/src/Messages"
)

type IHandler interface {
	HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message)
}

func HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	var handler IHandler
	switch msg.Type {
	case "preinit":
		handler = &PreInitHandler{}
	case "natify":
		handler = &NatifyHandler{}
	case "connect_ack":
		handler = &ConnectAckHandler{}
	case "ert_ack":
		handler = &ERTAckHandler{}
	case "init":
		handler = &InitHandler{}
	case "report":
		handler = &ReportHandler{}
	}
	if handler == nil {
		log.Printf("unhandled message: %s\n", msg.Type)
		return
	}
	handler.HandleMessage(core, outboundHandler, msg)
}
