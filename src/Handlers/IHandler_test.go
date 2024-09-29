package Handlers

import (
	"log"
	"net/netip"
	"testing"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type OutboundTestHandler struct {
	gotReportData      bool
	gotInitAck         bool
	gotDeadbeat        bool
	curretConnectIndex int
	connectClients     [2]*NatNegSessionClient
	connectAddresses   [2]netip.AddrPort
}

func (h *OutboundTestHandler) SendMessage(msg Messages.Message) {
	if msg.Type == "report_ack" {
		h.gotReportData = true
	} else if msg.Type == "init_ack" {
		h.gotInitAck = true
	} else if msg.Type == "ert" {
		/*var portType int = 0
		if msg.Message.(*Messages.ERTMessage).UnsolicitedPort {
			portType = 1
		}*/
		log.Printf("got ert req - %s %s\n", msg.DriverAddress, msg.Address)
	}
}
func (h *OutboundTestHandler) SendDeadbeatMessage(client *NatNegSessionClient) {
	h.gotDeadbeat = true
}
func (h *OutboundTestHandler) SendConnectMessage(client *NatNegSessionClient, ipAddress netip.AddrPort) {
	h.connectAddresses[h.curretConnectIndex] = ipAddress
	h.connectClients[h.curretConnectIndex] = client
	h.curretConnectIndex = h.curretConnectIndex + 1
}

func TestReport(t *testing.T) {
	var outboundHandler IOutboundHandler
	var testHandler *OutboundTestHandler = &OutboundTestHandler{}
	outboundHandler = testHandler

	var core NatNegCore
	core.Init(outboundHandler, 2, "11.11.11.11:1111", "11.11.11.11:2222", "22.22.22.22:3333")

	var msg Messages.Message
	var reportMsg Messages.ReportMessage
	reportMsg.Gamename = "gamename"
	msg.Message = &reportMsg
	msg.Type = "report"
	msg.Version = 4

	HandleMessage(core, outboundHandler, msg)

	if !testHandler.gotReportData {
		t.Errorf("report ack not sent")
	}
}

func TestInit_ExpectConnect_WithRetry_VerifyDeleteAfterAck(t *testing.T) {
	var outboundHandler IOutboundHandler
	var testHandler OutboundTestHandler = OutboundTestHandler{}
	outboundHandler = &testHandler

	var core NatNegCore
	core.Init(outboundHandler, 5, "11.11.11.11:1111", "11.11.11.11:2222", "22.22.22.22:3333")

	var cookie = 123321

	var msg Messages.Message
	var initMsg Messages.InitMessage
	initMsg.ClientIndex = 0
	initMsg.PrivateAddress = "10.1.1.1:7777"
	initMsg.PortType = 0
	initMsg.UseGamePort = 1

	msg.Message = &initMsg
	msg.Cookie = cookie
	msg.Type = "init"
	msg.Version = 2
	msg.DriverAddress = "127.0.0.1:11111"
	msg.Address = "25.25.25.25:6500"
	msg.Gamename = "test"
	HandleMessage(core, outboundHandler, msg)

	msg.DriverAddress = "127.0.0.2:11111"
	msg.Address = "25.25.25.25:6500"
	initMsg.PortType = 1
	msg.Message = &initMsg
	HandleMessage(core, outboundHandler, msg)

	msg.DriverAddress = "127.0.0.3:11111"
	msg.Address = "25.25.25.25:6500"
	initMsg.PortType = 2
	msg.Message = &initMsg
	HandleMessage(core, outboundHandler, msg)

	initMsg.ClientIndex = 1
	initMsg.PrivateAddress = "10.1.1.1:7777"
	initMsg.PortType = 0
	initMsg.UseGamePort = 1

	msg.Message = &initMsg
	msg.Type = "init"
	msg.Version = 2
	msg.DriverAddress = "127.0.0.1:11111"
	msg.Address = "66.25.25.25:6500"
	msg.Gamename = "test"
	HandleMessage(core, outboundHandler, msg)

	msg.DriverAddress = "127.0.0.2:11111"
	msg.Address = "66.25.25.25:6500"
	initMsg.PortType = 1
	msg.Message = &initMsg
	HandleMessage(core, outboundHandler, msg)

	msg.DriverAddress = "127.0.0.3:11111"
	msg.Address = "66.25.25.25:6500"
	initMsg.PortType = 2
	msg.Message = &initMsg
	HandleMessage(core, outboundHandler, msg)

	for i := 0; i < 10 && testHandler.curretConnectIndex != 2; i++ {
		core.Tick()
		time.Sleep(1 * time.Second)
	}

	if testHandler.curretConnectIndex != 2 {
		t.Errorf("Unexpected connect index")
		return
	}
	var gpInitAddr = testHandler.connectClients[0].findAddressInfoOfType(NN_SERVER_GP)
	if gpInitAddr.Address != testHandler.connectAddresses[1] {
		t.Errorf("Unexpected connect address: %s != %s", gpInitAddr.Address.String(), testHandler.connectAddresses[1].String())
	}

	gpInitAddr = testHandler.connectClients[1].findAddressInfoOfType(NN_SERVER_GP)
	if gpInitAddr.Address != testHandler.connectAddresses[0] {
		t.Errorf("Unexpected connect address: %s != %s", gpInitAddr.Address.String(), testHandler.connectAddresses[0].String())
	}

	//now wait for retry before sending ack back
	testHandler.curretConnectIndex = 0 //reset idx

	for i := 0; i < 10 && testHandler.curretConnectIndex != 2; i++ {
		core.Tick()
		time.Sleep(1 * time.Second)
	}

	if testHandler.curretConnectIndex != 2 {
		t.Errorf("Unexpected connect index - retry not sent")
		return
	}

	gpInitAddr = testHandler.connectClients[0].findAddressInfoOfType(NN_SERVER_GP)
	if gpInitAddr.Address != testHandler.connectAddresses[1] {
		t.Errorf("Unexpected connect address: %s != %s", gpInitAddr.Address.String(), testHandler.connectAddresses[1].String())
	}

	gpInitAddr = testHandler.connectClients[1].findAddressInfoOfType(NN_SERVER_GP)
	if gpInitAddr.Address != testHandler.connectAddresses[0] {
		t.Errorf("Unexpected connect address: %s != %s", gpInitAddr.Address.String(), testHandler.connectAddresses[0].String())
	}

	//send acks from both clients
	var connectAckMsg Messages.Message
	connectAckMsg.Cookie = cookie
	connectAckMsg.DriverAddress = "127.0.0.2:11111"
	connectAckMsg.Address = "66.25.25.25:6500"
	connectAckMsg.Type = "connect_ack"
	HandleMessage(core, outboundHandler, connectAckMsg)

	connectAckMsg.Address = "25.25.25.25:6500"
	HandleMessage(core, outboundHandler, connectAckMsg)

	var sess = core.findSessionByCookie(cookie)
	if sess != nil {
		t.Errorf("Session not deleted")
	}

	if !testHandler.gotInitAck {
		t.Errorf("didn't get init ack")
	}
}

func TestInit_ExpectDeadbeat(t *testing.T) {
	var outboundHandler IOutboundHandler
	var testHandler *OutboundTestHandler = &OutboundTestHandler{}
	outboundHandler = testHandler

	var core NatNegCore
	core.Init(outboundHandler, 2, "11.11.11.11:1111", "11.11.11.11:2222", "22.22.22.22:3333")

	var msg Messages.Message
	var initMsg Messages.InitMessage
	initMsg.ClientIndex = 1
	initMsg.PrivateAddress = "10.1.1.1:7777"
	initMsg.PortType = 0
	initMsg.UseGamePort = 1

	msg.Message = &initMsg
	msg.Type = "init"
	msg.Version = 2
	msg.DriverAddress = "127.0.0.1:11111"
	msg.Address = "25.25.25.25:6500"
	msg.Gamename = "test"
	HandleMessage(core, outboundHandler, msg)

	msg.DriverAddress = "127.0.0.2:11111"
	msg.Address = "25.25.25.25:6500"
	initMsg.PortType = 1
	msg.Message = &initMsg
	HandleMessage(core, outboundHandler, msg)

	msg.DriverAddress = "127.0.0.3:11111"
	msg.Address = "25.25.25.25:6500"
	initMsg.PortType = 2
	msg.Message = &initMsg
	HandleMessage(core, outboundHandler, msg)

	for i := 0; i < 10 && !testHandler.gotDeadbeat; i++ {
		core.Tick()
		time.Sleep(1 * time.Second)
	}

	if !testHandler.gotDeadbeat {
		t.Errorf("didn't get deadbeat")
	}
}
