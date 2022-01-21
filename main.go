package main

import (
	"encoding/json"
	"fmt"

	"github.com/ivahaev/amigo"
)

// Structs
type PeerStatus struct {
	Peer string `json:"peer"`
	PeerStatus string `json:"peer_status"`
	Address string `json:"peer_address"`
}

type Queue struct {
	CallerIDNum string `json:"caller_num,omitempty"`
	CallerIDName string `json:"caller_name,omitempty"`
	Count string `json:"queue_total_channels,omitempty"`
	Postition string `json:"queue_postition,omitempty"`
	Queue string `json:"queue,omitempty"`
	HoldTime string `json:"hold_time,omitempty"`
	RingTime string `json:"ring_time,omitempty"`
	TalkTime string `json:"talk_time,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
	AgentNumber string `json:"agent_number,omitempty"`
	Reason string `json:"end_reason,omitempty"`
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

// Global variables
var (
	call = Calls{}
)

// Helper function remove channel from Calls
func removeOutboundChannel(call *Calls, i int) {
	call.Outbound[i] = call.Outbound[len(call.Outbound)-1]
	call.Outbound = call.Outbound[:len(call.Outbound)-1]
}

// Helper function remove channel from Calls
func removeInboundChannel(call *Calls, i int) {
	call.Inbound[i] = call.Inbound[len(call.Inbound)-1]
	call.Inbound = call.Inbound[:len(call.Inbound)-1]
}

// Helper function to output json
func printJson (call interface{}) {
	jsonOutput, err := json.MarshalIndent(call, "", "   ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(jsonOutput))
}

// Creating hanlder functions
func DeviceStateChangeHandler(m map[string]string) {
	fmt.Printf("DeviceStateChange event received: %v\n", m)
}

func DefaultHandler(m map[string]string) {
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

func DialEventHandler(m map[string]string) {
	fmt.Printf("Dial Event received: %v\n", m)
}

func HangupEventHandler(m map[string]string) {
	// fmt.Printf("Hangup Event received: %v\n", m)
	
	// Outbound calls
	for i, v := range call.Outbound {
		if v.UID == m["Uniqueid"] {
			fmt.Printf("HANGUP %s \n",m["Cause-txt"])
			removeOutboundChannel(&call, i)
			printJson(call)
		}
	}

	// Inbound Calls
	for i, v := range call.Inbound {
		if v.UID == m["Uniqueid"] {
			fmt.Printf("HANGUP %s \n",m["Cause-txt"])
			removeInboundChannel(&call, i)
			printJson(call)
		}
	}

}

func QueueJoinEventHandler(m map[string]string) {
	for i, v := range call.Inbound {
		if  v.UID == m["Uniqueid"] {
			fmt.Printf("Q JOIN")
			call.Inbound[i].Queue.CallerIDName = m["CallerIDName"]
			call.Inbound[i].Queue.CallerIDNum = m["CallerIDNum"]
			call.Inbound[i].Queue.Count = m["Count"]
			call.Inbound[i].Queue.Postition = m["Position"]
			call.Inbound[i].Queue.Queue = m["Queue"]

			printJson(call)
		}
	}
}

func AgentConnectEventHandler(m map[string]string) {
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

func AgentCompleteEventHandler(m map[string]string) {
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

//TODO: QueueCallAbandoned Event

func NewStateEventHandler(m map[string]string) {

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

func NewChannelHandler(m map[string]string) {
	// fmt.Printf("New Channel Event: %v \n", m)
	// Outbound call
	if (m["Context"] == "from-internal" && m["Exten"] != "") {
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
	if (m["Context"] == "from-trunk" && m["Exten"] != "") {
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

func PeerStatusHandler(m map[string]string) {
	// fmt.Printf("PeerStatus event received: %v\n", m)

	peer := PeerStatus{
		Peer: m["Peer"],
		PeerStatus: m["PeerStatus"],
		Address: m["Address"],
	}

	pj, _ := json.Marshal(peer)

	fmt.Printf(string(pj)+"\n")


}

func main() {

	settings := &amigo.Settings{Username: "phpari", Password: "phpari", Host: "172.16.223.250"}
	a := amigo.New(settings)

	a.Connect()

	// Listen for connection events
	a.On("connect", func(message string) {
		fmt.Println("Connected", message)
	})
	a.On("error", func(message string) {
		fmt.Println("Connection error:", message)
	})


	// Registering handler function for event "DeviceStateChange"
	// a.RegisterHandler("DeviceStateChange", DeviceStateChangeHandler)
	//  a.RegisterHandler("PeerStatus", PeerStatusHandler)
	// a.RegisterHandler("Dial",DialEventHandler)

	a.RegisterHandler("Hangup",HangupEventHandler)
	a.RegisterHandler("Newchannel", NewChannelHandler)
	a.RegisterHandler("Newstate",NewStateEventHandler)


	// // Registering default handler function for all events.
	// a.RegisterDefaultHandler(DefaultHandler)
	a.RegisterHandler("Join", QueueJoinEventHandler)
	a.RegisterHandler("AgentConnect",AgentConnectEventHandler)
	a.RegisterHandler("AgentComplete",AgentCompleteEventHandler)


	// // Check if connected with Asterisk, will send Action "QueueSummary"
	// if a.Connected() {
	// 	result, err := a.Action(map[string]string{"Action": "QueueSummary", "ActionID": "Init"})
	// 	// If not error, processing result. Response on Action will follow in defined events.
	// 	// You need to catch them in event channel, DefaultHandler or specified HandlerFunction
	// 	fmt.Println(result, err)
	// }
	
	// do not exit main 
	forever := make(chan bool)
	<-forever
}