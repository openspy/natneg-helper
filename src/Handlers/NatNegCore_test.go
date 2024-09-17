package Handlers

import (
	"log"
	"net/netip"
	"testing"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type NNCoreTestOBH struct {
	gotDeadbeat       bool
	gotConnect        bool
	connectAddressIdx int
	connectAddress    [2]netip.AddrPort
}

func (c *NNCoreTestOBH) SendMessage(msg Messages.Message) {

}
func (c *NNCoreTestOBH) SendDeadbeatMessage(client *NatNegSessionClient) {
	//log.Printf("send deadbeat to: %s\n", client.Address)
	c.gotDeadbeat = true
}
func (c *NNCoreTestOBH) SendConnectMessage(client *NatNegSessionClient, ipAddress netip.AddrPort) {
	c.gotConnect = true
	c.connectAddress[c.connectAddressIdx] = ipAddress
	c.connectAddressIdx = c.connectAddressIdx + 1
}

func setup() (NatNegCore, *NNCoreTestOBH) {
	var obh *NNCoreTestOBH = &NNCoreTestOBH{}
	var core NatNegCore
	core.Init(obh, 2)
	return core, obh
}
func TestInit_GotPeers_OpenNATAll(t *testing.T) {
	core, obh := setup()

	var msg Messages.Message
	msg.Version = 2
	msg.Cookie = 111
	msg.Type = "init"
	msg.DriverAddress = "10.1.1.1:6666"
	msg.FromAddress = "127.0.0.1:7777"

	var initMsg Messages.InitMessage
	initMsg.LocalIP = 111
	initMsg.LocalPort = 7777

	//CLIENT 1
	initMsg.ClientIndex = 0
	initMsg.PortType = 0
	initMsg.UseGamePort = 1
	msg.Message = initMsg
	core.HandleInitMessage(msg) //NN1 / GamePort init - conn 1

	initMsg.PortType = 1
	msg.Message = initMsg
	msg.FromAddress = "127.0.0.1:12312"
	core.HandleInitMessage(msg) //NN1 / init1 init - conn 1

	msg.DriverAddress = "10.1.1.1:6667"
	msg.FromAddress = "127.0.0.1:7778"
	initMsg.PortType = 2
	msg.Message = initMsg
	core.HandleInitMessage(msg) //NN2 / init2 init - conn 1

	//CLIENT 2

	initMsg.ClientIndex = 1
	initMsg.PortType = 0
	msg.Message = initMsg
	msg.DriverAddress = "10.1.1.1:6666"
	msg.FromAddress = "25.25.25.25:7777"
	core.HandleInitMessage(msg) //NN1 / GamePort - conn 2

	msg.FromAddress = "25.25.25.25:22312"
	initMsg.PortType = 1
	msg.Message = initMsg
	core.HandleInitMessage(msg) //NN1 / init 1 - conn 2

	msg.DriverAddress = "10.1.1.1:6667"
	msg.FromAddress = "25.25.25.25:7778"
	initMsg.PortType = 2
	msg.Message = initMsg
	core.HandleInitMessage(msg) //NN1 / init 1 - conn 2

	if !obh.gotConnect {
		t.Errorf("Didn't get connect message")
	} else {
		log.Printf("got connect address 1: %s\n", obh.connectAddress[0].String())
		log.Printf("got connect address 2: %s\n", obh.connectAddress[1].String())
	}
}
func TestDeadbeat(t *testing.T) {
	core, obh := setup()

	var msg Messages.Message
	msg.Cookie = 111
	msg.Type = "init"
	msg.DriverAddress = "10.1.1.1:6666"
	msg.FromAddress = "127.0.0.1:7777"

	var initMsg Messages.InitMessage
	initMsg.LocalIP = 111
	initMsg.LocalPort = 666
	initMsg.UseGamePort = 1

	msg.Message = initMsg

	core.HandleInitMessage(msg)
	for i := range 60 {
		core.Tick()
		time.Sleep(1 * time.Second)
		if obh.gotDeadbeat {
			break
		}
	}
	if !obh.gotDeadbeat {
		t.Errorf("Didn't get deadbeat message")
	}

}
