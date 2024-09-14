package Handlers

import (
	"log"
	"net/netip"
)

type NATType int
type NATPromiscuity int
type NATMappingScheme int

const (
	NAT_TYPE_UNKNOWN NATType = iota
	NAT_TYPE_NO_NAT
	NAT_TYPE_FIREWALL_ONLY
	NAT_TYPE_FULL_CONE
	NAT_TYPE_RESTRICTED_CONE
	NAT_TYPE_PORT_RESTRICTED_CONE
	NAT_TYPE_SYMMETRIC
)

const (
	NAT_PROMISCUITY_PROMISCUOUS NATPromiscuity = iota
	NAT_PROMISCUITY_NOT_PROMISCUOUS
	NAT_PROMISCUITY_PORT_PROMISCUOUS
	NAT_PROMISCUITY_IP_PROMISCUOUS
	NAT_PROMISCUITY_PROMISCUITY_NOT_AVAILABLE
)
const (
	NAT_MAPPING_SCHEME_UNRECOGNIZED NATMappingScheme = iota
	NAT_MAPPING_SCHEME_PRIVATE_AS_PUBLIC
	NAT_MAPPING_SCHEME_CONSISTENT_PORT
	NAT_MAPPING_SCHEME_INCREMENTAL
	NAT_MAPPING_SCHEME_MIXED
)

type NatNegResolver struct {
}

func (c *NatNegResolver) portsMatch(expected uint16, addresses ...*NatNegSessionAddressInfo) bool {
	for i := range addresses {
		if addresses[i] != nil || addresses[i].Address.Port() != expected {
			return false
		}
	}
	return true
}
func (c *NatNegResolver) detectNAT_Version2(session NatNegSessionClient) (NATType, NATPromiscuity, NATMappingScheme) {
	var natType NATType = NAT_TYPE_UNKNOWN
	var promiscuity NATPromiscuity = NAT_PROMISCUITY_PROMISCUITY_NOT_AVAILABLE
	var natMappingScheme NATMappingScheme = NAT_MAPPING_SCHEME_UNRECOGNIZED

	var solicitedReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_NN1)
	var solicitedReply2 *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_NN2)
	var unsolicitedIPReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_UNSOLICITED_IP_PROBE)
	var unsolicitedIPPortReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_UNSOLICITED_IPPORT_PROBE)
	var unsolicitedPortReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_UNSOLICITED_PORT_PROBE)

	//var restricted bool = solicitedReply == nil
	var ipRestricted bool = unsolicitedIPReply == nil
	var portRestricted bool = unsolicitedPortReply == nil && unsolicitedIPPortReply == nil
	var publicAddrIsPrivateAddr bool = session.PrivateAddress.Addr() == solicitedReply.Address.Addr()

	/*var diff int = 0

	if unsolicitedIPPortReply != nil {
		diff = int(math.Abs(float64(solicitedReply.Address.Port() - unsolicitedIPPortReply.Address.Port())))
	}*/

	if !portRestricted && !ipRestricted && publicAddrIsPrivateAddr {
		natType = NAT_TYPE_NO_NAT
	} else if publicAddrIsPrivateAddr {
		natType = NAT_TYPE_FIREWALL_ONLY
	} else {
		// What type of NAT is it?
		var isSymmetric bool = solicitedReply.Address.Addr() != solicitedReply2.Address.Addr()
		if isSymmetric {
			natType = NAT_TYPE_SYMMETRIC
		} else if portRestricted {
			natType = NAT_TYPE_PORT_RESTRICTED_CONE
		} else if ipRestricted && !portRestricted {
			natType = NAT_TYPE_RESTRICTED_CONE
		} else if !ipRestricted && !portRestricted {
			natType = NAT_TYPE_FULL_CONE
		}
	}

	// What is the port mapping behavior?
	if c.portsMatch(session.PrivateAddress.Port(), solicitedReply, solicitedReply2, unsolicitedIPPortReply, unsolicitedIPReply, unsolicitedPortReply) {
		// Using private port as the public port.
		natMappingScheme = NAT_MAPPING_SCHEME_PRIVATE_AS_PUBLIC
	} else if c.portsMatch(solicitedReply.Address.Port(), solicitedReply2, unsolicitedIPPortReply, unsolicitedIPReply, unsolicitedPortReply) {
		// Using the same public port for all requests from the same private port.
		natMappingScheme = NAT_MAPPING_SCHEME_CONSISTENT_PORT
	} else if c.portsMatch(session.PrivateAddress.Port(), solicitedReply) && c.portsMatch(session.PrivateAddress.Port()+1, solicitedReply2) {
		// Using private port as the public port for the first mapping.
		// Using an incremental (+1) port mapping scheme there after.
		natMappingScheme = NAT_MAPPING_SCHEME_MIXED
	}

	return natType, promiscuity, natMappingScheme
}
func (c NatNegResolver) resolveNAT(natType NATType, natPromiscuity NATPromiscuity, mappingScheme NATMappingScheme) netip.AddrPort {
	v, _ := netip.ParseAddrPort("85.5.5.5:5555")
	return v
}
func (c *NatNegResolver) ResolveNAT(session NatNegSessionClient) netip.AddrPort {

	switch session.Version {
	case 2:
		natType, promiscuity, natMappingScheme := c.detectNAT_Version2(session)
		return c.resolveNAT(natType, promiscuity, natMappingScheme)
	default:
		log.Panicf("Unhandled NAT version")
	}
	return netip.AddrPort{}
}
