package Handlers

import (
	"openspy.net/natneg-helper/src/Messages"
)

/*
looks like pre-init is meant to wait for both clients to appear and possibly more initalization (matchup), but for now just say they are ready
*/

const (
	NN_PREINIT_WAITING_FOR_CLIENT  int = 0
	NN_PREINIT_WAITING_FOR_MATCHUP     = 1
	NN_PREINIT_READY                   = 2
)

type PreInitHandler struct {
}

func (b *PreInitHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	var preinit *Messages.PreInitMessage = msg.Message.(*Messages.PreInitMessage)
	preinit.State = NN_PREINIT_READY

	msg.Type = "preinit_ack"
	msg.Message = preinit
	outboundHandler.SendMessage(msg)
}
