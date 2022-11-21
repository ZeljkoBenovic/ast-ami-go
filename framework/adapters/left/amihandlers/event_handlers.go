package amihandlers

import (
	"os"
)

func (a *Adapter) newChannelHandler() {
	if err := a.amigo.RegisterHandler("Newchannel", func(m map[string]string) {
		// Outbound call
		if m["Context"] == a.config.OutboundContext && m["Exten"] != "" && m["Exten"] != "s" {
			a.amiEvents.Outbound[CallUID(m["Uniqueid"])] = OutboundCall{
				Type:         "OUTBOUND",
				CallerIDNum:  m["CallerIDNum"],
				CallerIDName: m["CallerIDName"],
				Context:      m["Context"],
				Exten:        m["Exten"],
				UID:          m["Uniqueid"],
				Event:        "NEW_OUTBOUND_CALL",
				EventCode:    NewOutboundCall,
				Timestamp:    convertTimeToUnixTime(m["TimeReceived"]),
			}

			a.logger.Debug("Call registered", "direction", "outbound", "event", "NEW_OUTBOUND_CALL", "data", m)
			a.logger.Debug("Events map", "event", "NEW_OUTBOUND_CALL", "map", a.amiEvents)
			a.logger.Info("Call registered",
				"event", "NEW_OUTBOUND_CALL",
				"direction", "outbound",
				"caller_id", m["CallerIDNum"],
				"call_id", m["Uniqueid"])

			a.sendDataToWebhook(m["Uniqueid"], outbound)
		}
		// Inbound call
		if m["Context"] == a.config.InboundContext && m["Exten"] != "" && m["Exten"] != "s" {
			a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = InboundCall{
				Type:         "INBOUND",
				CallerIDNum:  m["CallerIDNum"],
				CallerIDName: m["CallerIDName"],
				Context:      m["Context"],
				Exten:        m["Exten"],
				UID:          m["Uniqueid"],
				Event:        "NEW_INBOUND_CALL",
				EventCode:    NewInboundCall,
				Timestamp:    convertTimeToUnixTime(m["TimeReceived"]),
			}

			a.logger.Debug("Call registered", "event", "NEW_INBOUND_CALL", "direction", "inbound", "data", m)
			a.logger.Debug("Events map", "event", "NEW_INBOUND_CALL", "map", a.amiEvents)
			a.logger.Info("Call registered", "event", "NEW_INBOUND_CALL", "direction", "inbound",
				"caller_id", m["CallerIDNum"],
				"call_id", m["Uniqueid"])

			a.sendDataToWebhook(m["Uniqueid"], inbound)
		}
	}); err != nil {
		a.logger.Error("Could not register handler", "handler", "NEWCHANNEL")
		os.Exit(1)
	}
}

func (a *Adapter) hangupHandler() {
	if err := a.amigo.RegisterHandler("Hangup", func(m map[string]string) {
		if elem, ok := a.amiEvents.Outbound[CallUID(m["Uniqueid"])]; ok {
			// change call status
			elem.Event = "OUTBOUND_CALL_END"
			elem.EventCode = EndOutboundCall
			elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
			a.amiEvents.Outbound[CallUID(m["Uniqueid"])] = elem
			// send data
			a.sendDataToWebhook(m["Uniqueid"], outbound)
			// and delete element
			delete(a.amiEvents.Outbound, CallUID(m["Uniqueid"]))

			a.logger.Info("Call removed",
				"event", "OUTBOUND_CALL_END",
				"direction", "outbound",
				"caller_id", m["CallerIDNum"],
				"call_id", m["Uniqueid"])
			a.logger.Debug("Hangup event registered", "event", "OUTBOUND_CALL_END", "direction", "outbound", "data", m)
			a.logger.Debug("Call event deleted from map",
				"event", "OUTBOUND_CALL_END",
				"direction", "outbound",
				"call_id", m["Uniqueid"],
				"map", a.amiEvents)
		}

		if elem, ok := a.amiEvents.Inbound[CallUID(m["Uniqueid"])]; ok {
			// change call status
			elem.Event = "INBOUND_CALL_END"
			elem.EventCode = EndInboundCall
			elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
			a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = elem
			// send data
			a.sendDataToWebhook(m["Uniqueid"], inbound)
			// and delete element
			delete(a.amiEvents.Inbound, CallUID(m["Uniqueid"]))

			a.logger.Info("Call removed",
				"event", "INBOUND_CALL_END",
				"direction", "inbound",
				"caller_id", m["CallerIDNum"],
				"call_id", m["Uniqueid"])
			a.logger.Debug("Hangup event registered", "event", "INBOUND_CALL_END", "direction", "inbound", "event", m)
			a.logger.Debug("Call event deleted from map",
				"event", "INBOUND_CALL_END",
				"direction", "inbound",
				"call_id", m["Uniqueid"],
				"map", a.amiEvents)
		}
	}); err != nil {
		a.logger.Error("Could not register handler", "handler", "HANGUP")
		os.Exit(1)
	}
}

//nolint:dupl
func (a *Adapter) newStateHandler() {
	if err := a.amigo.RegisterHandler("Newstate", func(m map[string]string) {
		if elem, ok := a.amiEvents.Outbound[CallUID(m["Uniqueid"])]; ok {
			switch m["ChannelState"] {
			case "4":
				elem.Event = "RINGING"
				elem.EventCode = Ringing
				elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
				a.amiEvents.Outbound[CallUID(m["Uniqueid"])] = elem

				a.logger.Debug("Call state changed", "event", "RINGING",
					"direction", "outbound", "event", m)
				a.logger.Debug("Events map", "event", "RINGING",
					"direction", "outbound", "map", a.amiEvents)
				a.logger.Info("Call state changed", "event", "RINGING",
					"direction", "outbound",
					"caller_id", m["CallerIDNum"],
					"called_num", m["Exten"],
					"call_id", m["Uniqueid"])

				a.sendDataToWebhook(m["Uniqueid"], outbound)

			case "6":
				elem.Event = "ANSWERED"
				elem.EventCode = Answered
				elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
				elem.CallerIDName = m["CallerIDName"]
				a.amiEvents.Outbound[CallUID(m["Uniqueid"])] = elem

				a.logger.Debug("Call state changed", "event", "ANSWERED",
					"direction", "outbound", "event", m)
				a.logger.Debug("Events map", "event", "ANSWERED",
					"direction", "outbound", "map", a.amiEvents)
				a.logger.Info("Call state changed", "event", "ANSWERED",
					"direction", "outbound",
					"caller_id", m["CallerIDNum"],
					"called_num", m["Exten"],
					"call_id", m["Uniqueid"])

				a.sendDataToWebhook(m["Uniqueid"], outbound)
			}
		}

		if elem, ok := a.amiEvents.Inbound[CallUID(m["Uniqueid"])]; ok {
			switch m["ChannelState"] {
			case "4":
				elem.Event = "RINGING"
				elem.EventCode = Ringing
				elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
				a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = elem

				a.logger.Debug("Call state changed", "event", "RINGING",
					"direction", "inbound", "event", m)
				a.logger.Debug("Events map", "event", "RINGING",
					"direction", "inbound", "map", a.amiEvents)
				a.logger.Info("Call state changed", "event", "RINGING",
					"direction", "inbound",
					"caller_id", m["CallerIDNum"],
					"called_num", m["Exten"],
					"call_id", m["Uniqueid"])

				a.sendDataToWebhook(m["Uniqueid"], inbound)

			case "6":
				elem.Event = "ANSWERED"
				elem.EventCode = Answered
				elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
				elem.CallerIDName = m["CallerIDName"]
				a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = elem

				a.logger.Debug("Call state changed", "event", "ANSWERED",
					"direction", "inbound", "event", m)
				a.logger.Debug("Events map", "event", "ANSWERED",
					"direction", "inbound", "map", a.amiEvents)
				a.logger.Info("Call state changed", "event", "ANSWERED",
					"direction", "inbound",
					"caller_id", m["CallerIDNum"],
					"called_num", m["Exten"],
					"call_id", m["Uniqueid"])

				a.sendDataToWebhook(m["Uniqueid"], inbound)
			}
		}
	}); err != nil {
		a.logger.Error("Could not register handler", "handler", "NEWSTATE")
		os.Exit(1)
	}
}

func (a *Adapter) queueJoinEvent() {
	if err := a.amigo.RegisterHandler("Join", func(m map[string]string) {
		if elem, ok := a.amiEvents.Inbound[CallUID(m["Uniqueid"])]; ok {
			elem.Event = "QUEUE_JOIN"
			elem.EventCode = QueueJoin
			elem.Queue.CallerIDName = m["CallerIDName"]
			elem.Queue.CallerIDNum = m["CallerIDNum"]
			elem.Queue.Count = m["Count"]
			elem.Queue.Position = m["Position"]
			elem.Queue.Queue = m["Queue"]
			elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
			a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = elem

			a.logger.Debug("Call state changed", "event", "QUEUE_JOIN",
				"direction", "inbound", "event", m)
			a.logger.Debug("Events map", "event", "QUEUE_JOIN",
				"direction", "inbound", "map", a.amiEvents)
			a.logger.Info("Call state changed", "event", "QUEUE_JOIN",
				"direction", "inbound",
				"caller_id", m["CallerIDNum"],
				"queue_num", m["Queue"],
				"call_id", m["Uniqueid"])

			a.sendDataToWebhook(m["Uniqueid"], inbound)
		}
	}); err != nil {
		a.logger.Error("Could not register handler", "handler", "JOIN")
		os.Exit(1)
	}
}

func (a *Adapter) agentConnectEvent() {
	if err := a.amigo.RegisterHandler("AgentConnect", func(m map[string]string) {
		if elem, ok := a.amiEvents.Inbound[CallUID(m["Uniqueid"])]; ok {
			elem.Queue.HoldTime = m["HoldTime"]
			elem.Queue.RingTime = m["RingTime"]
			elem.Queue.AgentName = parseAgentName(m["MemberName"], a.logger)
			elem.Event = "AGENT_CONNECT"
			elem.EventCode = AgentConnect
			elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
			a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = elem

			a.logger.Debug("Call state changed", "event", "AGENT_CONNECT",
				"direction", "inbound", "event", m)
			a.logger.Debug("Events map", "event", "AGENT_CONNECT",
				"direction", "inbound", "map", a.amiEvents)
			a.logger.Info("Call state changed", "event", "AGENT_CONNECT",
				"direction", "inbound",
				"caller_id", m["CallerIDNum"],
				"agent_name", m["MemberName"],
				"call_id", m["Uniqueid"])

			a.sendDataToWebhook(m["Uniqueid"], inbound)
		}
	}); err != nil {
		a.logger.Error("Could not register handler", "handler", "AGENTCONNECT")
		os.Exit(1)
	}
}

func (a *Adapter) agentComplete() {
	if err := a.amigo.RegisterHandler("AgentComplete", func(m map[string]string) {
		if elem, ok := a.amiEvents.Inbound[CallUID(m["Uniqueid"])]; ok {
			elem.Queue.HoldTime = m["HoldTime"]
			elem.Queue.AgentName = parseAgentName(m["MemberName"], a.logger)
			elem.Queue.Reason = m["Reason"]
			elem.Queue.TalkTime = m["TalkTime"]
			elem.Event = "AGENT_COMPLETE"
			elem.EventCode = AgentComplete
			elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
			a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = elem

			a.logger.Debug("Call state changed", "event", "AGENT_COMPLETE",
				"direction", "inbound", "event", m)
			a.logger.Debug("Events map", "event", "AGENT_COMPLETE",
				"direction", "inbound", "map", a.amiEvents)
			a.logger.Info("Call state changed", "event", "AGENT_COMPLETE",
				"direction", "inbound",
				"caller_id", m["CallerIDNum"],
				"call_id", m["Uniqueid"])

			a.sendDataToWebhook(m["Uniqueid"], inbound)
		}
	}); err != nil {
		a.logger.Error("Could not register handler", "handler", "AGENTCOMPLETE")
		os.Exit(1)
	}
}

func (a *Adapter) queueAbandon() {
	if err := a.amigo.RegisterHandler("QueueCallerAbandon", func(m map[string]string) {
		if elem, ok := a.amiEvents.Inbound[CallUID(m["Uniqueid"])]; ok {
			elem.Queue.HoldTime = m["HoldTime"]
			elem.Queue.OriginalPosition = m["OriginalPosition"]
			elem.Queue.Position = m["Position"]
			elem.Queue.Queue = m["Queue"]
			elem.Event = "QUEUE_ABANDON"
			elem.EventCode = QueueAbandon
			elem.Timestamp = convertTimeToUnixTime(m["TimeReceived"])
			a.amiEvents.Inbound[CallUID(m["Uniqueid"])] = elem

			a.logger.Debug("Call state changed", "event", "QUEUE_ABANDON",
				"direction", "inbound", "event", m)
			a.logger.Debug("Events map", "event", "QUEUE_ABANDON",
				"direction", "inbound", "map", a.amiEvents)
			a.logger.Info("Call state changed", "event", "QUEUE_ABANDON",
				"direction", "inbound",
				"caller_id", m["CallerIDNum"],
				"hold_time", m["HoldTime"],
				"call_id", m["Uniqueid"])

			a.sendDataToWebhook(m["Uniqueid"], inbound)
		}
	}); err != nil {
		a.logger.Error("Could not register handler", "handler", "QUEUECALLERABANDON")
		os.Exit(1)
	}
}
