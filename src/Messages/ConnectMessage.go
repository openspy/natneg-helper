package Messages

const (
	FINISHED_NOERROR                    int = 0
	FINISHED_ERROR_DEAD_PARTNER             = 1
	FINISHED_ERROR_INIT_PACKETS_TIMEOUT     = 2
)

type ConnectMessage struct {
	RemoteIP    int `json:"remote_ip"`
	RemotePort  int `json:"remote_port"`
	GotYourData int `json:"got_your_data"`
	Finished    int `json:"finished"`
}
