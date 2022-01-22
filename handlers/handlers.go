package handlers

import "fmt"

//lint:ignore U1000 Ignore unused function as we use it for debugging
func (call *Calls) defaultHandler(m map[string]string) {
	// fmt.Printf("Event received: %v\n", m)	
	for _, v := range call.Outbound {
		if m["Uniqueid"] == v.UID {
			fmt.Println(m)
		}
	}
	for _, v := range call.Inbound {
		if m["Uniqueid"] == v.UID {
			fmt.Println(m)
		}
	}
}


func (call *Calls) hangupEventHandler(m map[string]string) {
	// fmt.Printf("Hangup Event received: %v\n", m)
	
	// Outbound calls
	for i, v := range call.Outbound {
		if v.UID == m["Uniqueid"] {
			call.Outbound[i].Event = "OUTBOUND_CALL_END"
			call.logger(call.Outbound[i])

			call.removeOutboundChannel(i)
		}
	}

	// Inbound Calls
	for i, v := range call.Inbound {
		if v.UID == m["Uniqueid"] {
			call.Inbound[i].Event = "INBOUND_CALL_END"
			call.logger(call.Inbound[i])

			call.removeInboundChannel(i)
		}
	}

}

func (call *Calls) queueJoinEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			call.Inbound[i].Queue.CallerIDName = m["CallerIDName"]
			call.Inbound[i].Queue.CallerIDNum = m["CallerIDNum"]
			call.Inbound[i].Queue.Count = m["Count"]
			call.Inbound[i].Queue.Position = m["Position"]
			call.Inbound[i].Queue.Queue = m["Queue"]
			call.Inbound[i].Event = "QUEUE_JOIN"

			call.logger(call.Inbound[i])
		}
	}
}

func (call *Calls) agentConnectEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			call.Inbound[i].Queue.HoldTime = m["HoldTime"]
			call.Inbound[i].Queue.RingTime = m["RingTime"]
			call.Inbound[i].Queue.AgentName = m["MemberName"]
			call.Inbound[i].Event = "AGENT_CONNECT"

			call.logger(call.Inbound[i])
		}
	}
}

func (call *Calls) agentCompleteEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			call.Inbound[i].Queue.HoldTime = m["HoldTime"]
			call.Inbound[i].Queue.AgentName = m["MemberName"]
			call.Inbound[i].Queue.Reason = m["Reason"]
			call.Inbound[i].Queue.TalkTime = m["TalkTime"]
			call.Inbound[i].Event = "AGENT_COMPLETE"

			call.logger(call.Inbound[i])
		}
	}
}

func (call *Calls) queueCallerAbandonEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			call.Inbound[i].Queue.HoldTime = m["HoldTime"]
			call.Inbound[i].Queue.OriginalPosition = m["OriginalPosition"]
			call.Inbound[i].Queue.Position = m["Position"]
			call.Inbound[i].Queue.Queue = m["Queue"]
			call.Inbound[i].Event = "QUEUE_ABANDON"

			call.logger(call.Inbound[i])
		}
	}
}

//TODO: QueueCallAbandoned Event

func (call *Calls) newStateEventHandler(m map[string]string) {

	// Outbound calls
	for i,v := range call.Outbound {
		if  v.UID == m["Uniqueid"] {
			switch m["ChannelState"] {
			case "4": // RINGING
				call.Outbound[i].Event = "RINGING"
				call.logger(call.Outbound[i])
			case "6": // ANSWERED
			call.Outbound[i].Event = "ANSWERED"
				call.logger(call.Outbound[i])
			}
		} 
	}

	// Inbound calls
	for i,v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			switch m["ChannelState"] {
			case "4": // RINGING
			call.Inbound[i].Event = "RINGING"
				call.logger(call.Inbound[i])
			case "6": // ANSWERED
			call.Inbound[i].Event = "ANSWERED"
				call.logger(call.Inbound[i])
			}
		} 
	}
}

func (call *Calls) newChannelHandler(m map[string]string) {
	// fmt.Printf("New Channel Event: %v \n", m)
	// Outbound call
	if (m["Context"] == call.OutboundContext && m["Exten"] != "") {
		newChannel := OutboundCall{
			CallerIDNum: m["CallerIDNum"],
			CallerIDName: m["CallerIDName"],
			Context: m["Context"],
			Exten: m["Exten"],
			UID: m["Uniqueid"],
			Event: "NEW_OUTBOUND_CALL",
		}
		call.logger(newChannel)
		call.Outbound = append(call.Outbound, newChannel)
	}

	// Inbound call
	if (m["Context"] == call.InboundContext && m["Exten"] != "") {
		newChannel := InboundCall{
			CallerIDNum: m["CallerIDNum"],
			CallerIDName: m["CallerIDName"],
			Context: m["Context"],
			Exten: m["Exten"],
			UID: m["Uniqueid"],
			Event: "NEW_INBOUND_CALL",
		}
		call.logger(newChannel)
		call.Inbound = append(call.Inbound, newChannel)
	}
}


