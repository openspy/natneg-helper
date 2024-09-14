package Handlers

import (
	"fmt"

	"openspy.net/natneg-helper/src/Messages"
)

type IHandler interface {
	HandleMessage(outboundHandler IOutboundHandler, msg Messages.Message)
}

func HandleMessage(outboundHandler IOutboundHandler, msg Messages.Message) {
	var handler IHandler
	fmt.Println("Handle message")
	switch msg.Type {
	case "preinit":
	case "connect":
	case "natify":
		handler = nil
	case "init":
		handler = &InitHandler{}
	case "report":
		handler = &ReportHandler{}
	}
	handler.HandleMessage(outboundHandler, msg)
}
