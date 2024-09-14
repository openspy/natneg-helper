package Handlers

import (
	"openspy.net/natneg-helper/src/Messages"
)

type ReportHandler struct {
	Version int         `json:"version"`
	Type    int         `json:"type"`
	Message interface{} `json:"message"`
}

func (b *ReportHandler) HandleMessage(outboundHandler IOutboundHandler, msg Messages.Message) {
	//var reportMsg Messages.ReportMessage = msg.Message.(Messages.ReportMessage)
	//fmt.Printf("aa %s\n", reportMsg.Gamename)

	msg.Type = "report_ack"
	outboundHandler.SendMessage(msg)
}
