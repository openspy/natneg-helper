package Handlers

import (
	"fmt"
	"log"
	"net/netip"
	"time"

	"openspy.net/natneg-helper/src/Messages"
)

type NNServerType int

const (
	NN_SERVER_GP                       NNServerType = iota //probe from natneg1 on gameport?? - maybe not needed
	NN_SERVER_NN1                                          //probe from natneg1 (matching ip/port)
	NN_SERVER_NN2                                          //probe from natneg2 (matching ip/port)
	NN_SERVER_NN3                                          //probe from natneg3 (matching ip/port)
	NN_SERVER_UNSOLICITED_IPPORT_PROBE                     //probe from NNCore (different port)
	NN_SERVER_UNSOLICITED_IP_PROBE                         //prove from NNCore (same port)
	NN_SERVER_UNSOLICITED_PORT_PROBE                       //probe from NN1
)

type NatNegSessionAddressInfo struct {
	RecvTime      time.Time
	ServerType    NNServerType
	Address       netip.AddrPort
	DriverAddress netip.AddrPort
}

type NatNegSessionClient struct {
	Version        int
	GotClient      bool
	UseGamePort    bool
	InitTime       time.Time
	PrivateAddress netip.AddrPort
	InitAddresses  []NatNegSessionAddressInfo

	//from CONNECT_PING unless ERT is available, then its ERT
	ConnectPingAddresses []NatNegSessionAddressInfo

	ConnectAddress  netip.AddrPort //final resolved address - for retries
	LastSentConnect time.Time
	ConnectAckTime  time.Time
}

func (c *NatNegSessionClient) findAddressInfoOfType(servType NNServerType) *NatNegSessionAddressInfo {
	switch servType {
	case NN_SERVER_GP:
		fallthrough
	case NN_SERVER_NN1:
		fallthrough
	case NN_SERVER_NN2:
		fallthrough
	case NN_SERVER_NN3:
		fallthrough
	default:
		for idx, _ := range c.InitAddresses {
			if c.InitAddresses[idx].ServerType == servType {
				return &c.InitAddresses[idx]
			}
		}
	case NN_SERVER_UNSOLICITED_IPPORT_PROBE:
		fallthrough
	case NN_SERVER_UNSOLICITED_IP_PROBE:
		fallthrough
	case NN_SERVER_UNSOLICITED_PORT_PROBE:
		for idx, _ := range c.ConnectPingAddresses {
			if c.ConnectPingAddresses[idx].ServerType == servType {
				return &c.ConnectPingAddresses[idx]
			}
		}
	}

	return nil
}
func (c *NatNegSessionClient) IsComplete() bool {
	var numExpected int = 2
	if c.Version >= 3 {
		numExpected = 3
	}
	if c.UseGamePort {
		numExpected = numExpected + 1
	}

	if len(c.InitAddresses) == numExpected {
		return true
	}
	return false
}

type NatNegSession struct {
	Version           int
	Cookie            int
	SessionCreateTime time.Time
	SessionClients    [2]NatNegSessionClient
	Resolver          INatNegResolver
}

func (s *NatNegSession) IsComplete() bool {
	return s.SessionClients[0].IsComplete() && s.SessionClients[1].IsComplete()
}

type NatNegCore struct {
	timer               *time.Ticker
	Sessions            map[int]*NatNegSession
	outboundHandler     IOutboundHandler
	resolver            INatNegResolver
	deadbeatTimeoutSecs int
	connectRetrySecs    int
}

func (c *NatNegCore) Init(obh IOutboundHandler, deadbeatTimeoutSecs int) {
	c.timer = time.NewTicker(time.Second)
	c.Sessions = make(map[int]*NatNegSession)
	c.outboundHandler = obh
	c.deadbeatTimeoutSecs = deadbeatTimeoutSecs
	c.connectRetrySecs = 5
	c.resolver = &NatNegResolver{}
}
func (c *NatNegCore) checkDeadbeats(currentTime time.Time) {
	for _, session := range c.Sessions {
		diff := currentTime.Sub(session.SessionCreateTime).Seconds()
		if diff > float64(c.deadbeatTimeoutSecs) {
			c.deleteSession(session.Cookie, true)
		}
	}
}

func (c *NatNegCore) createSession(msg Messages.Message) *NatNegSession {
	var session = &NatNegSession{}
	session.Cookie = msg.Cookie
	session.SessionCreateTime = time.Now()
	session.Version = msg.Version

	session.SessionClients[0].InitAddresses = make([]NatNegSessionAddressInfo, 0)
	session.SessionClients[1].InitAddresses = make([]NatNegSessionAddressInfo, 0)

	c.Sessions[msg.Cookie] = session
	return c.Sessions[msg.Cookie]
}

func (c *NatNegCore) Tick() {
	t := <-c.timer.C
	c.checkDeadbeats(t)
	c.checkConnectRetries(t)
}

func (c *NatNegCore) findSessionByCookie(cookie int) *NatNegSession {
	val, ok := c.Sessions[cookie]
	if ok {
		return val
	}
	return nil
}

func (c *NatNegCore) sendDeadbeat(victim *NatNegSessionClient) {
	c.outboundHandler.SendDeadbeatMessage(victim)
}
func (c *NatNegCore) deleteSession(cookie int, sendDeadbeat bool) {
	var session = c.Sessions[cookie]
	delete(c.Sessions, cookie)

	if sendDeadbeat {
		for _, clientSession := range session.SessionClients {
			if clientSession.GotClient {
				c.sendDeadbeat(&clientSession)
			}
		}
	}

}

func (c *NatNegCore) HandleInitTimeout() {

}
func (c *NatNegCore) HandleInitMessage(msg Messages.Message) {
	var session = c.findSessionByCookie(msg.Cookie)
	if session == nil {
		session = c.createSession(msg)
	}

	var initMsg = msg.Message.(Messages.InitMessage)
	var clientSession = &session.SessionClients[initMsg.ClientIndex]
	clientSession.InitTime = time.Now()
	clientSession.GotClient = true
	clientSession.Version = msg.Version

	if initMsg.UseGamePort != 0 {
		clientSession.UseGamePort = true
	}

	ipport, parseerr := netip.ParseAddrPort(msg.FromAddress)
	if parseerr != nil {
		log.Panicf("Failed to parse IP Port: %s\n", parseerr.Error())
	}
	//clientSession.Address = ipport

	var info NatNegSessionAddressInfo
	info.Address = ipport
	info.RecvTime = time.Now()
	switch initMsg.PortType {
	case 0:
		info.ServerType = NN_SERVER_GP
	case 1:
		info.ServerType = NN_SERVER_NN1
	case 2:
		info.ServerType = NN_SERVER_NN2
	case 3:
		info.ServerType = NN_SERVER_NN3
	}
	clientSession.InitAddresses = append(clientSession.InitAddresses, info)

	clientSession.PrivateAddress = netip.MustParseAddrPort("192.168.11.22:5511")

	fmt.Printf("got cookie: %d, idx: %d, addr: %s, type: %d, private: %s\n", msg.Cookie, initMsg.ClientIndex, ipport.String(), initMsg.PortType, clientSession.PrivateAddress.String())

	if session.IsComplete() {
		c.sendNegotiatedConnection(session)
	}

}

func (c *NatNegCore) checkConnectRetries(currentTime time.Time) {
	for _, session := range c.Sessions {
		diff := currentTime.Sub(session.SessionClients[0].LastSentConnect).Seconds()
		if session.SessionClients[0].ConnectAckTime.IsZero() && diff > float64(c.connectRetrySecs) {
			session.SessionClients[0].LastSentConnect = time.Now()
			c.outboundHandler.SendConnectMessage(&session.SessionClients[0], session.SessionClients[0].ConnectAddress)
		}

		diff = currentTime.Sub(session.SessionClients[1].LastSentConnect).Seconds()
		if session.SessionClients[1].ConnectAckTime.IsZero() && diff > float64(c.connectRetrySecs) {
			session.SessionClients[1].LastSentConnect = time.Now()
			c.outboundHandler.SendConnectMessage(&session.SessionClients[1], session.SessionClients[1].ConnectAddress)
		}
	}
}

func (c *NatNegCore) sendNegotiatedConnection(session *NatNegSession) {
	var resolved = c.resolver.ResolveNAT(session.SessionClients[0])
	session.SessionClients[1].ConnectAddress = resolved
	fmt.Printf("Send conn message 1: %s\n", resolved.String())
	session.SessionClients[1].LastSentConnect = time.Now()
	c.outboundHandler.SendConnectMessage(&session.SessionClients[1], resolved)

	resolved = c.resolver.ResolveNAT(session.SessionClients[1])
	session.SessionClients[0].ConnectAddress = resolved
	fmt.Printf("Send conn message 2: %s\n", resolved.String())
	session.SessionClients[0].LastSentConnect = time.Now()
	c.outboundHandler.SendConnectMessage(&session.SessionClients[0], resolved)

}
