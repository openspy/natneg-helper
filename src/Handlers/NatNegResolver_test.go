package Handlers

import (
	"net/netip"
	"testing"
	"time"
)

func TestResolver_ExpectedNoNAT(t *testing.T) {
	var session NatNegSessionClient
	session.Version = 2
	session.UseGamePort = true
	session.GotClient = true
	session.PrivateAddress = netip.MustParseAddrPort("25.25.25.25:6500")

	var initItem NatNegSessionAddressInfo
	initItem.Address = netip.MustParseAddrPort("25.25.25.25:6500")
	initItem.RecvTime = time.Now()

	initItem.ServerType = NN_SERVER_GP
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.1:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN1
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.2:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN2
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.3:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_IPPORT_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_IP_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_PORT_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	var resolver NatNegResolver
	natType, _, _ := resolver.detectNAT_Version2(session)

	if natType != NAT_TYPE_NO_NAT {
		t.Errorf("Got unexpected name type")
	}
}

func TestResolver_ExpecteFirewallOnly(t *testing.T) {
	var session NatNegSessionClient
	session.Version = 2
	session.UseGamePort = true
	session.GotClient = true
	session.PrivateAddress = netip.MustParseAddrPort("25.25.25.25:6500")

	var initItem NatNegSessionAddressInfo
	initItem.Address = netip.MustParseAddrPort("25.25.25.25:6500")
	initItem.RecvTime = time.Now()

	initItem.ServerType = NN_SERVER_GP
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.1:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN1
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.2:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN2
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.3:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	var resolver NatNegResolver
	natType, _, _ := resolver.detectNAT_Version2(session)

	if natType != NAT_TYPE_FIREWALL_ONLY {
		t.Errorf("Got unexpected name type")
	}
}

func TestResolver_ExpecteFullCone(t *testing.T) {
	var session NatNegSessionClient
	session.Version = 2
	session.UseGamePort = true
	session.GotClient = true
	session.PrivateAddress = netip.MustParseAddrPort("192.168.10.55:6500")

	var initItem NatNegSessionAddressInfo
	initItem.Address = netip.MustParseAddrPort("25.25.25.25:6500")
	initItem.RecvTime = time.Now()

	initItem.ServerType = NN_SERVER_GP
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.1:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN1
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.2:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN2
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.3:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_IPPORT_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_IP_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_PORT_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	var resolver NatNegResolver
	natType, _, _ := resolver.detectNAT_Version2(session)

	if natType != NAT_TYPE_FULL_CONE {
		t.Errorf("Got unexpected name type")
	}
}

func TestResolver_ExpecteSymmetric(t *testing.T) {
	var session NatNegSessionClient
	session.Version = 2
	session.UseGamePort = true
	session.GotClient = true
	session.PrivateAddress = netip.MustParseAddrPort("192.168.10.55:6500")

	var initItem NatNegSessionAddressInfo
	initItem.Address = netip.MustParseAddrPort("25.25.25.25:6500")
	initItem.RecvTime = time.Now()

	initItem.ServerType = NN_SERVER_GP
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.1:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN1
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.2:11111")
	initItem.Address = netip.MustParseAddrPort("65.25.25.25:6500")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN2
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.3:11111")
	initItem.Address = netip.MustParseAddrPort("5.25.25.25:6500")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	initItem.Address = netip.MustParseAddrPort("185.25.25.25:6500")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	initItem.Address = netip.MustParseAddrPort("195.25.25.25:6500")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_IPPORT_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_IP_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_PORT_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	var resolver NatNegResolver
	natType, _, _ := resolver.detectNAT_Version2(session)

	if natType != NAT_TYPE_SYMMETRIC {
		t.Errorf("Got unexpected name type")
	}
}

func TestResolver_ExpecteRestrictedPortCone(t *testing.T) {
	var session NatNegSessionClient
	session.Version = 2
	session.UseGamePort = true
	session.GotClient = true
	session.PrivateAddress = netip.MustParseAddrPort("192.168.10.55:6500")

	var initItem NatNegSessionAddressInfo
	initItem.Address = netip.MustParseAddrPort("25.25.25.25:6500")
	initItem.RecvTime = time.Now()

	initItem.ServerType = NN_SERVER_GP
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.1:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN1
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.2:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN2
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.3:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_IP_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	var resolver NatNegResolver
	natType, _, _ := resolver.detectNAT_Version2(session)

	if natType != NAT_TYPE_PORT_RESTRICTED_CONE {
		t.Errorf("Got unexpected name type")
	}
}

func TestResolver_ExpecteRestrictedCone(t *testing.T) {
	var session NatNegSessionClient
	session.Version = 2
	session.UseGamePort = true
	session.GotClient = true
	session.PrivateAddress = netip.MustParseAddrPort("192.168.10.55:6500")

	var initItem NatNegSessionAddressInfo
	initItem.Address = netip.MustParseAddrPort("25.25.25.25:6500")
	initItem.RecvTime = time.Now()

	initItem.ServerType = NN_SERVER_GP
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.1:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN1
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.2:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN2
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.3:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_NN3
	initItem.DriverAddress = netip.MustParseAddrPort("127.0.0.4:11111")
	session.InitAddresses = append(session.InitAddresses, initItem)
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	initItem.ServerType = NN_SERVER_UNSOLICITED_PORT_PROBE
	initItem.DriverAddress = netip.MustParseAddrPort("172.16.26.26:11111")
	session.ERTAddresses = append(session.ERTAddresses, initItem)

	var resolver NatNegResolver
	natType, _, _ := resolver.detectNAT_Version2(session)

	if natType != NAT_TYPE_RESTRICTED_CONE {
		t.Errorf("Got unexpected name type")
	}
}
