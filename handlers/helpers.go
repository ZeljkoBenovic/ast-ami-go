package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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

func logger (event interface{}) {
	jsonOutput, _ := json.Marshal(event)
	log.Println(string(jsonOutput))
}

//TODO: Add logging to file

//lint:ignore U1000 Ignore unused function as we use it for debugging
func printJson (call interface{}) {
	jsonOutput, err := json.MarshalIndent(call, "", "   ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(jsonOutput))
}

func debugPrint(msg map[string]string) {
	fmt.Println(msg)
}