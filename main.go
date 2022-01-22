package main

import (
	"fmt"

	"github.com/ZeljkoBenovic/ast-ami-go/handlers"
	"github.com/ivahaev/amigo"
)


func main() {

	// init AMI connection
	settings := &amigo.Settings{Username: "phpari", Password: "phpari", Host: "172.16.223.250"}
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
	callHandlers.RegisterHandlers(a)
	
	// do not exit main 
	forever := make(chan bool)
	<-forever
}