package Handlers

import (
	"log"
	"net/netip"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type ERTAckHandler struct {
}

func (b *ERTAckHandler) HandleMessage(core NatNegCore, outboundHandler IOutboundHandler, msg Messages.Message) {

	var nnType NNServerType = NN_SERVER_UNSOLICITED_PORT_PROBE

	/*portProbeDriver := os.Getenv("UNSOLICITED_PORT_PROBE_DRIVER")
	ipProbeDriver := os.Getenv("UNSOLICITED_IP_PROBE_DRIVER")
	ipPortProbeDriver := os.Getenv("UNSOLICITED_IPPORT_PROBE_DRIVER")*/

	portProbeDriver, ipProbeDriver, ipPortProbeDriver := core.GetERTDrivers()

	if msg.DriverAddress == ipPortProbeDriver {
		nnType = NN_SERVER_UNSOLICITED_IPPORT_PROBE
	} else if msg.DriverAddress == ipProbeDriver {
		nnType = NN_SERVER_UNSOLICITED_IP_PROBE
	} else if msg.DriverAddress == portProbeDriver {
		nnType = NN_SERVER_UNSOLICITED_PORT_PROBE
	}

	addr, addrErr := netip.ParseAddrPort(msg.Address)

	log.Printf("[%s] ERT ACK - type: %d\n", msg.Address, nnType)

	if addrErr != nil {
		log.Printf("ERTAckHandler got invalid address: %s\n", addrErr.Error())
		return
	}

	var session = core.findSessionByCookie(msg.Cookie)
	if session == nil {
		log.Printf("Got ERT request for invalid cookie %d - from %s\n", msg.Cookie, msg.Address)
		return
	}

	var ertRecord NatNegSessionAddressInfo
	ertRecord.ServerType = nnType
	ertRecord.Address = addr
	ertRecord.RecvTime = time.Now()

	driverAddr, driverAddrErr := netip.ParseAddrPort(msg.DriverAddress)

	if driverAddrErr != nil {
		log.Printf("ERTAckHandler got invalid address: %s - from %s\n", driverAddrErr.Error(), msg.Address)
		return
	}
	ertRecord.DriverAddress = driverAddr

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
