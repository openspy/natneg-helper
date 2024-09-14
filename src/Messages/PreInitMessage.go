package Messages

type PreInitMessage struct {
	ClientIndex int `json:"clientindex"`
	State       int `json:"state"`
	ClientID    int `json:"client_id"`
}
