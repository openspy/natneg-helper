package Messages

type ReportMessage struct {
	PortType         int    `json:"port_type"`
	ClientIndex      int    `json:"client_index"`
	NegResult        int    `json:"neg_result"`
	NATType          int    `json:"nat_type"`
	NATMappingScheme int    `json:"nat_mapping_scheme"`
	Gamename         string `json:"gamename"`
}
