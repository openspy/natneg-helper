package Messages

type InitMessage struct {
	PortType       int    `json:"porttype"`
	ClientIndex    int    `json:"clientindex"`
	UseGamePort    int    `json:"use_gameport"`
	PrivateAddress string `json:"private_address"`
}
