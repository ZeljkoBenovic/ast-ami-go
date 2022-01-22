package handlers


type Queue struct {
	CallerIDNum string `json:"caller_num,omitempty"`
	CallerIDName string `json:"caller_name,omitempty"`
	Count string `json:"queue_total_channels,omitempty"`
	Position string `json:"queue_position,omitempty"`
	Queue string `json:"queue,omitempty"`
	HoldTime string `json:"hold_time,omitempty"`
	RingTime string `json:"ring_time,omitempty"`
	TalkTime string `json:"talk_time,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
	AgentNumber string `json:"agent_number,omitempty"`
	Reason string `json:"end_reason,omitempty"`
	OriginalPosition string `json:"original_position,omitempty"`
}

type InboundCall struct {
	CallerIDNum string `json:"caller_id_num"`
	CallerIDName string `json:"caller_id_name"`
	Context string `json:"context"`
	Exten string `json:"did"`
	UID string `json:"uid"`
	Queue Queue `json:"queue,omitempty"`
}

type OutboundCall struct {
	CallerIDNum string `json:"extension"`
	CallerIDName string `json:"extension_name"`
	Context string `json:"context"`
	Exten string `json:"called_num"`
	UID string `json:"uid"`
}

type Calls struct {
	Outbound []OutboundCall `json:"outbound"`
	Inbound []InboundCall `json:"inbound"`
}