package Handlers

import (
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type ConnectAckHandler struct {
}

func (b *ConnectAckHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	var sess = core.findSessionByCookie(msg.Cookie)
	var client = core.GetClientFromMessage(msg)
	if client == nil {
		return
	}

	client.ConnectAckTime = time.Now()

	if !sess.SessionClients[0].ConnectAckTime.IsZero() && !sess.SessionClients[1].ConnectAckTime.IsZero() {
		core.deleteSession(msg.Cookie, false)
	}
}
