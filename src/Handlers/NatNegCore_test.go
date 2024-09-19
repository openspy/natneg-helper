package Handlers

import (
	"log"
	"net/netip"
	"testing"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

const REMOTE_DRIVER string = "66.66.66.66:6666"

type NNCoreTestOBH struct {
	gotDeadbeat       bool
	gotConnect        bool
	connectAddressIdx int
	connectAddress    [2]netip.AddrPort

	gotSolicatedERT         bool
	gotUnsolicitedPortERT   bool
	gotUnsolicitedIPPortERT bool

	core *NatNegCore
}

func (c *NNCoreTestOBH) SendMessage(msg Messages.Message) {
	if msg.Type == "ert" {
		var ertMsg *Messages.ERTMessage = msg.Message.(*Messages.ERTMessage)
		var unsolictedPort = ertMsg.UnsolicitedPort
		if msg.DriverAddress == REMOTE_DRIVER { //is remote / unsolicited
			if unsolictedPort {
				c.gotUnsolicitedIPPortERT = true
			}
		} else { //solicited
			c.gotSolicatedERT = true
			if unsolictedPort {
				c.gotUnsolicitedPortERT = true
			}
		}

	}
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

func setup(timeout int) (NatNegCore, *NNCoreTestOBH) {
	var obh *NNCoreTestOBH = &NNCoreTestOBH{}
	var core NatNegCore
	core.Init(obh, timeout, REMOTE_DRIVER)
	return core, obh
}
func TestInit_GotPeers_OpenNATAll(t *testing.T) {
	core, obh := setup(2)

	obh.core = &core

	var msg Messages.Message
	msg.Version = 2
	msg.Cookie = 111
	msg.Type = "init"
	msg.DriverAddress = "10.1.1.1:6666"
	msg.Address = "127.0.0.1:7777"

	var initMsg Messages.InitMessage
	initMsg.LocalIP = 0x19191919
	initMsg.LocalPort = 7777

	//CLIENT 1
	initMsg.ClientIndex = 0
	initMsg.PortType = 0
	initMsg.UseGamePort = 1
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN1 / GamePort init - conn 1

	initMsg.PortType = 1
	msg.Message = &initMsg
	msg.Address = "127.0.0.1:12312"
	core.HandleInitMessage(msg) //NN1 / init1 init - conn 1

	msg.DriverAddress = "10.1.1.1:6667"
	msg.Address = "127.0.0.1:7778"
	initMsg.PortType = 2
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN2 / init2 init - conn 1

	//CLIENT 2

	initMsg.ClientIndex = 1
	initMsg.PortType = 0
	msg.Message = &initMsg
	msg.DriverAddress = "10.1.1.1:6666"
	msg.Address = "25.25.25.25:7777"
	core.HandleInitMessage(msg) //NN1 / GamePort - conn 2

	msg.Address = "25.25.25.25:22312"
	initMsg.PortType = 1
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN1 / init 1 - conn 2

	msg.DriverAddress = "10.1.1.1:6667"
	msg.Address = "25.25.25.25:7778"
	initMsg.PortType = 2
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN1 / init 1 - conn 2

	core.Tick()

	if obh.gotConnect {
		t.Errorf("got connect too early")
	}

	if !obh.gotSolicatedERT || !obh.gotUnsolicitedIPPortERT || !obh.gotUnsolicitedPortERT {
		t.Errorf("didn't get ERT test")
	}

	var ertHandler ERTAckHandler

	//send unsolicited port - solicited IP response
	msg.Type = "ert_ack"
	var ertMsg Messages.ERTMessage
	ertMsg.UnsolicitedPort = true
	msg.Message = &ertMsg
	ertHandler.HandleMessage(core, obh, msg)

	//send unsolicted port - unsolicited IP response
	msg.DriverAddress = REMOTE_DRIVER
	ertHandler.HandleMessage(core, obh, msg)

	//send solicted port - unsolicited IP response
	msg.DriverAddress = REMOTE_DRIVER
	ertMsg.UnsolicitedPort = false
	msg.Message = &ertMsg
	ertHandler.HandleMessage(core, obh, msg)

	core.Tick()

	time.Sleep(5 * time.Second)
	core.Tick()

	if !obh.gotConnect {
		t.Errorf("Didn't get connect message")
	} else {
		log.Printf("got connect address 1: %s\n", obh.connectAddress[0].String())
		log.Printf("got connect address 2: %s\n", obh.connectAddress[1].String())
	}

	if obh.gotDeadbeat {
		t.Errorf("got unexpected deadbeat msg")
	}
}

func TestInit_GotPeers_xx(t *testing.T) {
	core, obh := setup(2)

	obh.core = &core

	var msg Messages.Message
	msg.Version = 2
	msg.Cookie = 111
	msg.Type = "init"
	msg.DriverAddress = "10.1.1.1:6666"
	msg.Address = "127.0.0.1:7777"

	var initMsg Messages.InitMessage
	initMsg.LocalIP = 0x19191919
	initMsg.LocalPort = 7777

	//CLIENT 1
	initMsg.ClientIndex = 0
	initMsg.PortType = 0
	initMsg.UseGamePort = 1
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN1 / GamePort init - conn 1

	initMsg.PortType = 1
	msg.Message = &initMsg
	msg.Address = "127.0.0.1:12312"
	core.HandleInitMessage(msg) //NN1 / init1 init - conn 1

	msg.DriverAddress = "10.1.1.1:6667"
	msg.Address = "127.0.0.1:7778"
	initMsg.PortType = 2
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN2 / init2 init - conn 1

	//CLIENT 2

	initMsg.ClientIndex = 1
	initMsg.PortType = 0
	msg.Message = &initMsg
	msg.DriverAddress = "10.1.1.1:6666"
	msg.Address = "25.25.25.25:7777"
	core.HandleInitMessage(msg) //NN1 / GamePort - conn 2

	msg.Address = "25.25.25.25:22312"
	initMsg.PortType = 1
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN1 / init 1 - conn 2

	msg.DriverAddress = "10.1.1.1:6667"
	msg.Address = "25.25.25.25:7778"
	initMsg.PortType = 2
	msg.Message = &initMsg
	core.HandleInitMessage(msg) //NN1 / init 1 - conn 2

	core.Tick()

	if obh.gotConnect {
		t.Errorf("got connect too early")
	}

	if !obh.gotSolicatedERT || !obh.gotUnsolicitedIPPortERT || !obh.gotUnsolicitedPortERT {
		t.Errorf("didn't get ERT test")
	}

	var ertHandler ERTAckHandler

	//send unsolicited port - solicited IP response
	msg.Type = "ert_ack"
	var ertMsg Messages.ERTMessage
	ertMsg.UnsolicitedPort = true
	msg.Message = &ertMsg
	ertHandler.HandleMessage(core, obh, msg)

	for i := 0; i < 120; i++ {
		core.Tick()
		time.Sleep(1 * time.Second)
		if obh.gotConnect {
			break
		}
	}

	if !obh.gotConnect {
		t.Errorf("Didn't get connect message")
	} else {
		log.Printf("got connect address 1: %s\n", obh.connectAddress[0].String())
		log.Printf("got connect address 2: %s\n", obh.connectAddress[1].String())
	}

	if obh.gotDeadbeat {
		t.Errorf("got unexpected deadbeat msg")
	}
}
func TestDeadbeat(t *testing.T) {
	core, obh := setup(2)

	var msg Messages.Message
	msg.Cookie = 111
	msg.Type = "init"
	msg.DriverAddress = "10.1.1.1:6666"
	msg.Address = "127.0.0.1:7777"

	var initMsg Messages.InitMessage
	initMsg.LocalIP = 111
	initMsg.LocalPort = 666
	initMsg.UseGamePort = 1

	msg.Message = &initMsg

	core.HandleInitMessage(msg)
	for i := 0; i < 60; i++ {
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

func TestNatifyReq(t *testing.T) {
	core, obh := setup(2)

	var msg Messages.Message
	msg.Cookie = 111
	msg.Type = "natify"
	msg.DriverAddress = "10.1.1.1:6666"
	msg.Address = "127.0.0.1:7777"

	var initMsg Messages.InitMessage
	initMsg.PortType = 0

	msg.Message = &initMsg

	var handler NatifyHandler
	handler.HandleMessage(core, obh, msg)

}
