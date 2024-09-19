package Handlers

import (
	"log"
	"net/netip"

	"openspy.net/natneg-helper/src/Messages"
)

const (
	NATIFY_COOKIE int = 777
)

type NatifyHandler struct {
}

func (b *NatifyHandler) sendRemoteERT(outboundHandler IOutboundHandler, driverAddress string, address netip.AddrPort, unsolicitedPort bool) {
	var msg Messages.Message
	msg.Cookie = NATIFY_COOKIE
	msg.DriverAddress = driverAddress

	msg.Address = address.String()

	var ertMsg Messages.InitMessage
	if unsolicitedPort {
		ertMsg.UseGamePort = 1
	}

	msg.Message = &ertMsg
	msg.Type = "ert"

	outboundHandler.SendMessage(msg)
}
func (b *NatifyHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	if msg.Cookie != NATIFY_COOKIE {
		log.Printf("Skipping natify message due to not using natify cookie")
		return
	}

	address, addressError := netip.ParseAddrPort(msg.Address)
	if addressError != nil {
		return
	}
	var natify = msg.Message.(*Messages.InitMessage)
	switch natify.PortType {
	case 0:
		fallthrough
	case 1:
		b.sendRemoteERT(outboundHandler, core.unsolictedPortERTDriver, address, true)
	case 2:
		b.sendRemoteERT(outboundHandler, core.unsolicteIPERTDriver, address, false)
	case 3:
		b.sendRemoteERT(outboundHandler, core.unsolictedIPPortERTDriver, address, true)
	}

}
