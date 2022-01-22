package handlers

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