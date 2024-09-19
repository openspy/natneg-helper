package Handlers

import (
	"net/netip"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type ERTAckHandler struct {
}

func (b *ERTAckHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {
	var unsolicitedIP bool = false
	if msg.DriverAddress == core.unsolictedERTDriver {
		unsolicitedIP = true
	}

	var nnType NNServerType = NN_SERVER_UNSOLICITED_PORT_PROBE
	var unsolictedPort bool = msg.Message.(Messages.ERTMessage).UnsolicitedPort

	if unsolicitedIP && unsolictedPort {
		nnType = NN_SERVER_UNSOLICITED_IPPORT_PROBE
	} else if unsolicitedIP {
		nnType = NN_SERVER_UNSOLICITED_IP_PROBE
	} else if unsolictedPort {
		nnType = NN_SERVER_UNSOLICITED_PORT_PROBE
	}

	var addr = netip.MustParseAddrPort(msg.Address)

	var session = core.findSessionByCookie(msg.Cookie)

	var ertRecord NatNegSessionAddressInfo
	ertRecord.ServerType = nnType
	ertRecord.Address = addr
	ertRecord.RecvTime = time.Now()
	ertRecord.DriverAddress = netip.MustParseAddrPort(msg.DriverAddress)
	if session != nil {
		if session.SessionClients[0].PublicIP == addr.Addr() {
			if session.SessionClients[0].findAddressInfoOfType(nnType) == nil {
				session.SessionClients[0].ERTAddresses = append(session.SessionClients[0].ERTAddresses, ertRecord)
			}

		} else if session.SessionClients[1].PublicIP == addr.Addr() {
			if session.SessionClients[1].findAddressInfoOfType(nnType) == nil {
				session.SessionClients[1].ERTAddresses = append(session.SessionClients[1].ERTAddresses, ertRecord)
			}
		}
	}

}
