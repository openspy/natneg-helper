package Handlers

import (
	"openspy.net/natneg-helper/src/Messages"
)

type InitHandler struct {
	Version int         `json:"version"`
	Type    int         `json:"type"`
	Message interface{} `json:"message"`
}

func (b *InitHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	msg.Type = "init_ack"
	outboundHandler.SendMessage(msg)

	core.HandleInitMessage(msg)

}
