package Handlers

import (
	"log"
	"net/netip"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type ConnectAckHandler struct {
}

func (b *ConnectAckHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	var sess = core.findSessionByCookie(msg.Cookie)
	var address = sess.SessionClients[0].findAddressInfoOfType(NN_SERVER_GP)
	fromAddr, _ := netip.ParseAddrPort(msg.FromAddress)
	log.Printf("0: %s == %s\n", fromAddr.Addr().String(), address.Address.Addr().String())
	if address.Address.Addr() == fromAddr.Addr() {
		sess.SessionClients[0].ConnectAckTime = time.Now()

	} else {
		address = sess.SessionClients[1].findAddressInfoOfType(NN_SERVER_GP)
		log.Printf("1: %s == %s\n", fromAddr.Addr().String(), address.Address.Addr().String())
		if address.Address.Addr() == fromAddr.Addr() {
			sess.SessionClients[1].ConnectAckTime = time.Now()
		}
	}

	if !sess.SessionClients[0].ConnectAckTime.IsZero() && !sess.SessionClients[1].ConnectAckTime.IsZero() {
		core.deleteSession(msg.Cookie, false)
	}
}
