package handlers

import (
	"log"

	"github.com/ivahaev/amigo"
)

func (call *Calls) RegisterHandlers(amigo *amigo.Amigo) {
	errors := make(map[string]error)
	// first poing of entry is new channel event
	errors["new_channel"] = amigo.RegisterHandler("Newchannel", call.newChannelHandler)
	errors["new_state"] = amigo.RegisterHandler("Newstate", call.newStateEventHandler)
	errors["hangup"] = amigo.RegisterHandler("Hangup", call.hangupEventHandler)
	errors["new_exten"] = amigo.RegisterHandler("NewExten", call.newExtenEventHandler)

	// // Registering default handler function for all events.
	//amigo.RegisterDefaultHandler(call.defaultHandler)

	// queue handlers
	errors["join"] = amigo.RegisterHandler("Join", call.queueJoinEventHandler)
	errors["agent_connect"] = amigo.RegisterHandler("AgentConnect", call.agentConnectEventHandler)
	errors["agent_complete"] = amigo.RegisterHandler("AgentComplete", call.agentCompleteEventHandler)
	errors["caller_abandon"] = amigo.RegisterHandler("QueueCallerAbandon", call.queueCallerAbandonEventHandler)

	for i, v := range errors {
		if v != nil {
			log.Fatalln("Could not register handler: ", i, "Error: ", v.Error())
		}
	}

	// Registering handler function for event "DeviceStateChange"
	// a.RegisterHandler("DeviceStateChange", DeviceStateChangeHandler)
	//  a.RegisterHandler("PeerStatus", PeerStatusHandler)
	// a.RegisterHandler("Dial",DialEventHandler)

	// // Check if connected with Asterisk, will send Action "QueueSummary"
	// if a.Connected() {
	// 	result, err := a.Action(map[string]string{"Action": "QueueSummary", "ActionID": "Init"})
	// 	// If not error, processing result. Response on Action will follow in defined events.
	// 	// You need to catch them in event channel, DefaultHandler or specified HandlerFunction
	// 	fmt.Println(result, err)
	// }
}
