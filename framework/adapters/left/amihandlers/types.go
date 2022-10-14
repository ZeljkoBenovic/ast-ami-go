package amihandlers

type Queue struct {
	CallerIDNum      string `json:"caller_num,omitempty"`
	CallerIDName     string `json:"caller_name,omitempty"`
	Count            string `json:"queue_total_channels,omitempty"`
	Position         string `json:"queue_position,omitempty"`
	Queue            string `json:"queue,omitempty"`
	HoldTime         string `json:"hold_time,omitempty"`
	RingTime         string `json:"ring_time,omitempty"`
	TalkTime         string `json:"talk_time,omitempty"`
	AgentName        string `json:"agent_name,omitempty"`
	AgentNumber      string `json:"agent_number,omitempty"`
	Reason           string `json:"end_reason,omitempty"`
	OriginalPosition string `json:"original_position,omitempty"`
}

type InboundCall struct {
	CallerIDNum  string    `json:"caller_id_num"`
	CallerIDName string    `json:"caller_id_name"`
	Context      string    `json:"-"`
	Exten        string    `json:"did"`
	UID          string    `json:"uid"`
	Queue        Queue     `json:"queue,omitempty"`
	Event        string    `json:"event"`
	EventCode    EventCode `json:"event_code"`
}

type OutboundCall struct {
	CallerIDNum  string    `json:"extension"`
	CallerIDName string    `json:"extension_name"`
	Context      string    `json:"-"`
	Exten        string    `json:"called_num"`
	UID          string    `json:"uid"`
	Event        string    `json:"event"`
	EventCode    EventCode `json:"event_code"`
}

type CallUID string

type EventCode int

const (
	NewInboundCall  EventCode = 1
	NewOutboundCall EventCode = 2
	Ringing         EventCode = 3
	Answered        EventCode = 4
	QueueJoin       EventCode = 5
	AgentConnect    EventCode = 6
	AgentComplete   EventCode = 7
	QueueAbandon    EventCode = 8
	EndInboundCall  EventCode = 9
	EndOutboundCall EventCode = 10
)

type Calls struct {
	Outbound map[CallUID]OutboundCall `json:"outbound"`
	Inbound  map[CallUID]InboundCall  `json:"inbound"`
}
