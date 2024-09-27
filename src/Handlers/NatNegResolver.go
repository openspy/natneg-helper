package Handlers

import (
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

func NATTypeToString(natType NATType) string {
	switch natType {
	case NAT_TYPE_UNKNOWN:
		return "UNKNOWN"
	case NAT_TYPE_NO_NAT:
		return "NO_NAT"
	case NAT_TYPE_FIREWALL_ONLY:
		return "FIREWALL_ONLY"
	case NAT_TYPE_FULL_CONE:
		return "FULL_CONE"
	case NAT_TYPE_RESTRICTED_CONE:
		return "RESTRICTED_CONE"
	case NAT_TYPE_PORT_RESTRICTED_CONE:
		return "PORT_RESTRICTED_CONE"
	case NAT_TYPE_SYMMETRIC:
		return "SYMMETRIC"
	}
	return ""
}

type NatNegResolver struct {
}

func (c *NatNegResolver) portsMatch(expected uint16, addresses ...*NatNegSessionAddressInfo) bool {
	for i := range addresses {
		if addresses[i] == nil || addresses[i].Address.Port() != expected {
			return false
		}
	}
	return true
}
func (c *NatNegResolver) detectNAT_Version2(session NatNegSessionClient) (NATType, NATPromiscuity, NATMappingScheme) {
	var natType NATType = NAT_TYPE_UNKNOWN
	var promiscuity NATPromiscuity = NAT_PROMISCUITY_PROMISCUITY_NOT_AVAILABLE
	var natMappingScheme NATMappingScheme = NAT_MAPPING_SCHEME_UNRECOGNIZED

	var solicitedReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_GP)
	if solicitedReply == nil {
		solicitedReply = session.findAddressInfoOfType(NN_SERVER_NN1)
	}
	var solicitedReply2 *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_NN2)
	var unsolicitedIPReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_UNSOLICITED_IP_PROBE)
	var unsolicitedIPPortReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_UNSOLICITED_IPPORT_PROBE)
	var unsolicitedPortReply *NatNegSessionAddressInfo = session.findAddressInfoOfType(NN_SERVER_UNSOLICITED_PORT_PROBE)

	//var restricted bool = solicitedReply == nil
	var ipRestricted bool = unsolicitedIPReply == nil
	var portRestricted bool = unsolicitedPortReply == nil && unsolicitedIPPortReply == nil
	var publicAddrIsPrivateAddr bool = solicitedReply != nil && session.PrivateAddress.IsValid() && session.PrivateAddress.Addr() == solicitedReply.Address.Addr()

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
		var isSymmetric bool = solicitedReply != nil && solicitedReply2 != nil && solicitedReply.Address.Addr() != solicitedReply2.Address.Addr()
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

func (c NatNegResolver) resolveNAT(session NatNegSessionClient, natType NATType, natPromiscuity NATPromiscuity, mappingScheme NATMappingScheme) netip.AddrPort {
	var serv = NN_SERVER_NN1
	if session.UseGamePort {
		serv = NN_SERVER_GP
	}
	var rootServer = session.findAddressInfoOfType(serv)
	if rootServer == nil {
		return netip.AddrPort{}
	}
	var unsolicitedAddress netip.AddrPort
	if session.PrivateAddress.Port() != 0 {
		unsolicitedAddress = netip.AddrPortFrom(rootServer.Address.Addr(), session.PrivateAddress.Port())
	} else {
		unsolicitedAddress = rootServer.Address
	}

	var returnAddress netip.AddrPort
	switch natType {
	case NAT_TYPE_NO_NAT:
		fallthrough
	case NAT_TYPE_FIREWALL_ONLY:
		fallthrough
	case NAT_TYPE_PORT_RESTRICTED_CONE:
		fallthrough
	case NAT_TYPE_RESTRICTED_CONE:
		fallthrough
	case NAT_TYPE_FULL_CONE:
		returnAddress = unsolicitedAddress

	//these nat types are unlikely to work
	default:
		fallthrough
	case NAT_TYPE_SYMMETRIC:
		fallthrough
	case NAT_TYPE_UNKNOWN:
		returnAddress = unsolicitedAddress
	}

	return returnAddress

}
func (c *NatNegResolver) DetectNAT(session NatNegSessionClient) (NATType, NATPromiscuity, NATMappingScheme) {
	switch session.Version {
	default:
		fallthrough
	case 2:
		return c.detectNAT_Version2(session)
		/*default:
		log.Panicf("Unhandled NAT version")*/
	}
	//return NAT_TYPE_UNKNOWN, NAT_PROMISCUITY_PROMISCUITY_NOT_AVAILABLE, NAT_MAPPING_SCHEME_UNRECOGNIZED
}
func (c *NatNegResolver) ResolveNAT(session NatNegSessionClient) netip.AddrPort {

	natType, promiscuity, natMappingScheme := c.DetectNAT(session)

	return c.resolveNAT(session, natType, promiscuity, natMappingScheme)
}
