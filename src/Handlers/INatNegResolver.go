package Handlers

import "net/netip"

type INatNegResolver interface {
	ResolveNAT(session NatNegSessionClient) netip.AddrPort
}
