package amihandlers

type Queue struct {
	CallerIDNum      string `json:"caller_num"`
	CallerIDName     string `json:"caller_name"`
	Count            string `json:"queue_total_channels"`
	Position         string `json:"queue_position"`
	Queue            string `json:"queue"`
	HoldTime         string `json:"hold_time"`
	RingTime         string `json:"ring_time"`
	TalkTime         string `json:"talk_time"`
	AgentName        string `json:"agent_name"`
	AgentNumber      string `json:"agent_number"`
	Reason           string `json:"end_reason"`
	OriginalPosition string `json:"original_position"`
}

type InboundCall struct {
	CallerIDNum  string    `json:"caller_id_num"`
	CallerIDName string    `json:"caller_id_name"`
	Context      string    `json:"-"`
	Exten        string    `json:"did"`
	UID          string    `json:"uid"`
	Event        string    `json:"event"`
	Recording    string    `json:"recording"`
	Timestamp    int64     `json:"timestamp"`
	EventCode    EventCode `json:"event_code"`
	Queue        Queue     `json:"queue"`
}

type OutboundCall struct {
	CallerIDNum  string    `json:"extension"`
	CallerIDName string    `json:"extension_name"`
	Context      string    `json:"-"`
	Exten        string    `json:"called_num"`
	UID          string    `json:"uid"`
	Event        string    `json:"event"`
	Recording    string    `json:"recording"`
	Timestamp    int64     `json:"timestamp"`
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
