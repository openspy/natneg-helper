package Handlers

import (
	"openspy.net/natneg-helper/src/Messages"
)

type ReportHandler struct {
}

func (b *ReportHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	//var reportMsg Messages.ReportMessage = msg.Message.(Messages.ReportMessage)
	//fmt.Printf("aa %s\n", reportMsg.Gamename)

	msg.Type = "report_ack"
	outboundHandler.SendMessage(msg)
}
