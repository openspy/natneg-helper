package Handlers

import (
	"openspy.net/natneg-helper/src/Messages"
)

type ReportHandler struct {
}

func (b *ReportHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	msg.Type = "report_ack"
	outboundHandler.SendMessage(msg)
}
