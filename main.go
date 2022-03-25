package main

import (
	"log"

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
	a.On("error", func(message string) {
		log.Fatalln("Connection error:", message)
	})
	
	// set handlers
	callHandlers := handlers.Calls{
		InboundContext: config.InboundContext, 
		OutboundContext: config.OutboundContext, 
		LogFileLocation: config.LogFileLocation,
		WebhookURL: config.WebhookURL,
		WebhookMethod: config.WebhookMethod,
		Debug: config.Debug,
	}
	callHandlers.RegisterHandlers(a)
	
	// do not exit main 
	forever := make(chan bool)
	<-forever
}