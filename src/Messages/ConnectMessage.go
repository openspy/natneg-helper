package Messages

type ConnectMessage struct {
	RemoteIP    int `json:"remote_ip"`
	RemotePort  int `json:"remote_port"`
	GotYourData int `json:"got_your_data"`
	Finished    int `json:"finished"`
}
