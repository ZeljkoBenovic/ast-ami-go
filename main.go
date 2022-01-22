package main

import (
	"fmt"

	"github.com/ZeljkoBenovic/ast-ami-go/config"
	"github.com/ZeljkoBenovic/ast-ami-go/handlers"
	"github.com/ivahaev/amigo"
)


func main() {

	// get configuration parameters
	config := config.GetConfig()
	// init AMI connection
	settings := &amigo.Settings{Username: config.Username, Password: config.Password, Host: config.Host, Port: config.Port}
	a := amigo.New(settings)
	
	// connect to AMI server
	a.Connect()
	
	// Listen for connection events
	a.On("connect", func(message string) {
		fmt.Println("Connected", message)
	})
	a.On("error", func(message string) {
		fmt.Println("Connection error:", message)
	})
	
	// set handlers
	callHandlers := handlers.Calls{}
	callHandlers.SetInboundContext("from-trunk")
	callHandlers.SetOutboundContext("from-internal")
	callHandlers.RegisterHandlers(a)
	
	// do not exit main 
	forever := make(chan bool)
	<-forever
}