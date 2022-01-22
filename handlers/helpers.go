package handlers

import (
	"encoding/json"
	"log"
	"os"
)

// Helper function remove channel from Calls
func (call *Calls) removeOutboundChannel(i int) {
	call.Outbound[i] = call.Outbound[len(call.Outbound)-1]
	call.Outbound = call.Outbound[:len(call.Outbound)-1]
}

// Helper function remove channel from Calls
func (call *Calls) removeInboundChannel(i int) {
	call.Inbound[i] = call.Inbound[len(call.Inbound)-1]
	call.Inbound = call.Inbound[:len(call.Inbound)-1]
}

func (call *Calls) logger (event interface{}) {
	var logFile *os.File
	var err error
	// initialze log file location
	if call.LogFileLocation == "" {
		logFile, err = os.OpenFile("call_event_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else {
		logFile, err = os.OpenFile(call.LogFileLocation, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	// set log file location
	log.SetOutput(logFile)

	// log content to file
	jsonOutput, _ := json.Marshal(event)
	log.Println(string(jsonOutput))
}