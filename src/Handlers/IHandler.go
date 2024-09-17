package Handlers

import (
	"fmt"

	"openspy.net/natneg-helper/src/Messages"
)

type IHandler interface {
	HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message)
}

func HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	var handler IHandler
	fmt.Println("Handle message")
	switch msg.Type {
	case "preinit":
	case "connect":
	case "natify":
		handler = nil
	case "connect_ack":
		handler = &ConnectAckHandler{}
	case "init":
		handler = &InitHandler{}
	case "report":
		handler = &ReportHandler{}
	}
	handler.HandleMessage(core, outboundHandler, msg)
}
