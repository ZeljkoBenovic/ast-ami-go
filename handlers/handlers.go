package handlers

import (
	"fmt"
)

//lint:ignore U1000 Ignore unused function as we use it for debugging
func (call *Calls) defaultHandler(m map[string]string) {
	// fmt.Printf("Event received: %v\n", m)	
	for _, v := range call.Outbound {
		if m["Uniqueid"] == v.UID {
			debugPrint(m)
		}
	}
	for _, v := range call.Inbound {
		if m["Uniqueid"] == v.UID {
			debugPrint(m)
		}
	}
}


func (call *Calls) hangupEventHandler(m map[string]string) {
	// fmt.Printf("Hangup Event received: %v\n", m)
	
	// Outbound calls
	for i, v := range call.Outbound {
		if v.UID == m["Uniqueid"] {
			fmt.Printf("HANGUP %s \n",m["Cause-txt"])
			call.removeOutboundChannel(i)
			printJson(call)
		}
	}

	// Inbound Calls
	for i, v := range call.Inbound {
		if v.UID == m["Uniqueid"] {
			fmt.Printf("HANGUP %s \n",m["Cause-txt"])
			call.removeInboundChannel(i)
			printJson(call)
		}
	}

}

func (call *Calls) queueJoinEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			fmt.Printf("Q JOIN")
			call.Inbound[i].Queue.CallerIDName = m["CallerIDName"]
			call.Inbound[i].Queue.CallerIDNum = m["CallerIDNum"]
			call.Inbound[i].Queue.Count = m["Count"]
			call.Inbound[i].Queue.Position = m["Position"]
			call.Inbound[i].Queue.Queue = m["Queue"]

			printJson(call)
		}
	}
}

func (call *Calls) agentConnectEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			fmt.Printf("Q AGENT CONNECT")
			call.Inbound[i].Queue.HoldTime = m["HoldTime"]
			call.Inbound[i].Queue.RingTime = m["RingTime"]
			call.Inbound[i].Queue.AgentName = m["MemberName"]

			printJson(call)
		}
	}
}

func (call *Calls) agentCompleteEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			fmt.Printf("Q AGENT COMPLETE")
			call.Inbound[i].Queue.HoldTime = m["HoldTime"]
			call.Inbound[i].Queue.AgentName = m["MemberName"]
			call.Inbound[i].Queue.Reason = m["Reason"]
			call.Inbound[i].Queue.TalkTime = m["TalkTime"]

			printJson(call)
		}
	}
}

func (call *Calls) queueCallerAbandonEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			fmt.Printf("Q CALL ABANDONED")
			call.Inbound[i].Queue.HoldTime = m["HoldTime"]
			call.Inbound[i].Queue.OriginalPosition = m["OriginalPosition"]
			call.Inbound[i].Queue.Position = m["Position"]
			call.Inbound[i].Queue.Queue = m["Queue"]

			printJson(call)
		}
	}
}

//TODO: QueueCallAbandoned Event

func (call *Calls) newStateEventHandler(m map[string]string) {

	// Outbound calls
	for _,v := range call.Outbound {
		if  v.UID == m["Uniqueid"] {
			// fmt.Println(m)
			switch m["ChannelState"] {
			case "4": // RINGING
				fmt.Println("RINGING")
				printJson(call)
				// fmt.Printf("Client %s is RINGING. Called by %s",v.Exten,v.CallerIDNum)
			case "6": // ANSWERED
				fmt.Println("ANSWERED")
				printJson(call)
				// fmt.Printf("Client %s ANSWERED. Called by %s",v.Exten,v.CallerIDNum)
			}
		} 
	}

	// Inbound calls
	for _,v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			// fmt.Println(m)
			switch m["ChannelState"] {
			case "4": // RINGING
				fmt.Println("RINGING")
				printJson(call)
				// fmt.Printf("Client %s is RINGING. Called by %s",v.Exten,v.CallerIDNum)
			case "6": // ANSWERED
				fmt.Println("ANSWERED")
				printJson(call)
				// fmt.Printf("Client %s ANSWERED. Called by %s",v.Exten,v.CallerIDNum)
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
		}
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
		}
		call.Inbound = append(call.Inbound, newChannel)
	}
}


