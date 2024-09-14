package Messages

type InitMessage struct {
	PortType    int `json:"type"`
	ClientIndex int `json:"client_index"`
	UseGamePort int `json:"use_gameport"`
	LocalIP     int `json:"localip"`
	LocalPort   int `json:"localport"`
}
