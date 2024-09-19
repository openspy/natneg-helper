package Messages

import (
	"encoding/json"
	"log"
)

type Message struct {
	Address       string      `json:"address"`
	DriverAddress string      `json:"driver_address"`
	Hostname      string      `json:"hostname"`
	Gamename      string      `json:"gamename"`
	Version       int         `json:"version"`
	Type          string      `json:"type"`
	Cookie        int         `json:"cookie"`
	Message       interface{} `json:"data"`
}

func (b *Message) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typ); err != nil {
		log.Printf("got err: %s\n", err)
		return err
	}

	switch typ.Type {
	case "natify":
		fallthrough
	case "ert":
		fallthrough
	case "ert_ack":
		fallthrough
	case "init":
		b.Message = new(InitMessage)
	case "connect":
		b.Message = new(ConnectMessage)
	case "preinit":
		b.Message = new(PreInitMessage)
	case "report":
		b.Message = new(ReportMessage)
	default:
		b.Message = nil
	}

	type tmp Message // avoids infinite recursion
	return json.Unmarshal(data, (*tmp)(b))
}
