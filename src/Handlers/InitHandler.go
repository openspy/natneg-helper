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
	//var reportMsg Messages.ReportMessage = msg.Message.(Messages.ReportMessage)
	//fmt.Printf("aa %s\n", reportMsg.Gamename)

	// select {
	// case <-time.After(2 * time.Second):
	// 	fmt.Printf("aaa\n")
	// 	/*default:
	// 	fmt.Printf("bbb\n")*/
	// }

	msg.Type = "init_ack"
	outboundHandler.SendMessage(msg)

	core.HandleInitMessage(msg)
}
