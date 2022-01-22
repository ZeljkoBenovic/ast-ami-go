package handlers

import "github.com/ivahaev/amigo"

func (call *Calls) RegisterHandlers(amigo *amigo.Amigo) {
	
	// first poing of entry is new channel event
	amigo.RegisterHandler("Newchannel", call.newChannelHandler)
	amigo.RegisterHandler("Newstate",call.newStateEventHandler)
	amigo.RegisterHandler("Hangup",call.hangupEventHandler)


	// // Registering default handler function for all events.
	// amigo.RegisterDefaultHandler(call.defaultHandler)

	// queue handlers
	amigo.RegisterHandler("Join", call.queueJoinEventHandler)
	amigo.RegisterHandler("AgentConnect",call.agentConnectEventHandler)
	amigo.RegisterHandler("AgentComplete",call.agentCompleteEventHandler)
	amigo.RegisterHandler("QueueCallerAbandon",call.queueCallerAbandonEventHandler)



	

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