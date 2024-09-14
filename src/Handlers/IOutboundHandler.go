package Handlers

import (
	"net/netip"

	"openspy.net/natneg-helper/src/Messages"
)

type IOutboundHandler interface {
	SendMessage(msg Messages.Message)
	SendDeadbeatMessage(client *NatNegSessionClient)
	SendConnectMessage(client *NatNegSessionClient, ipAddress netip.AddrPort)
}
