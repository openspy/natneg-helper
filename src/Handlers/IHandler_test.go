package Handlers

/*
test cases:
    * nn version 2 - nn init - no peer - expect deadbeat
    * nn version 2 - nn init - with peer - expect connection - send ack
*/

/*import (
	"testing"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type OutboundTestHandler struct {
	gotReportData bool
	gotInitAck    bool
}

func (h *OutboundTestHandler) SendMessage(msg Messages.Message) {
	if msg.Type == "report_ack" {
		h.gotReportData = true
	} else if msg.Type == "init_ack" {
		h.gotInitAck = true
	}

}

func TestReport(t *testing.T) {
	var outboundHandler IOutboundHandler
	var testHandler *OutboundTestHandler = &OutboundTestHandler{}
	outboundHandler = testHandler

	var msg Messages.Message
	var reportMsg Messages.ReportMessage
	reportMsg.Gamename = "gamename"
	msg.Message = reportMsg
	msg.Type = "report"
	msg.Version = 4

	HandleMessage(outboundHandler, msg)

	if !testHandler.gotReportData {
		t.Errorf("report ack not sent")
	}
}

func TestInit(t *testing.T) {
	var outboundHandler IOutboundHandler
	var testHandler *OutboundTestHandler = &OutboundTestHandler{}
	outboundHandler = testHandler

	var msg Messages.Message
	var initMsg Messages.InitMessage
	initMsg.ClientIndex = 1
	initMsg.LocalIP = 0x0A000001
	initMsg.LocalPort = 7777
	initMsg.PortType = 3
	initMsg.UseGamePort = 1

	msg.Message = initMsg
	msg.Type = "init"
	msg.Version = 4

	HandleMessage(outboundHandler, msg)

	time.Sleep(10 * time.Second)
	if !testHandler.gotInitAck {
		t.Errorf("init ack not sent")
	}
}
*/
