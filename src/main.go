package main

/*
Incoming messages:
	init:
		* use 3 NN cores to send CONNECT_PING
	connect_ping:
		* msg sent from NN core when user responds to connect ping on given udp ip/port
	connect_ack

Outgoing:
	connect_ping: send a connect_ping to a given IP
	connect
*/

/*
	Upon NN INIT:
		* send connect pings to all NN cores including from this service

		*  if version 4?? send ERT (for now from this service - ideally from a second service)

		*  when all responses come back from both peers (or 1/both peers timed out) //XXX: MOVE THIS??
			* determine who should be the primary peer to be connected
			* send connect/deadbeat messages
	Upon NN_PREINIT:
		unknown... gamespy SDK says it signals a queueing mechanism... maybe the server indicates when to begin
		send preinit ack - this contains 3 states
			NN_PREINIT_WAITING_FOR_CLIENT
			NN_PREINIT_WAITING_FOR_MATCHUP
				* both just tell the natneg client to wait
			NN_PREINIT_READY
				* tells natneg client to send init

		* for now just send preinit ready upon getting preinit
	Upon NN_REPORT:
		* send report ack
	Upon NN_CONNECT_PING:
		* use this info for natneg connect response
	Upon NN_CONNECT_ACK:
		* cancel connect resend
	Upon NN_NATIFY_REQUEST:
		* unknown // DiscoverReachability
*/

func main() {
	/*var outboundHandler AMQPOutboundHandler
	fmt.Println("Hello world")
	var handler Handlers.ReportHandler
	var msg Messages.Message
	var reportMsg Messages.ReportMessage
	reportMsg.Gamename = "hello"
	msg.Message = reportMsg
	msg.Type = "report"
	msg.Version = 4

	handler.HandleMessage(&outboundHandler, msg)*/
}
