package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func (call *Calls) logger (event interface{}) {
	var (
		logFile *os.File
		err error
	)

	jsonOutput, _ := json.Marshal(event)

	if call.Debug {
		fmt.Println("EVENT: "+string(jsonOutput))
	}

	// initialze log file location
	logFile, err = os.OpenFile(call.LogFileLocation, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(fmt.Errorf("could not setup log file: %v", err))
	}

	// write logs to log file
	log.SetOutput(logFile)
	// log event to file
	log.Println("EVENT: "+string(jsonOutput))

	// if webhook endpoint is defined send to it
	if call.WebhookURL != "" {
		// create new request endpoint
		req, err := http.NewRequest(call.WebhookMethod, call.WebhookURL, bytes.NewBuffer(jsonOutput))
		if err != nil {
			log.Println(fmt.Errorf("could not create new request endpoint: %w", err))
		}

		// set json header
		req.Header.Set("Content-Type", "application/json")
		// TODO: set optional authentication header - barer token

		// set new http client
		client := &http.Client{}
		// send the request and wait for responce
		resp, respErr := client.Do(req)
		if respErr != nil {
			log.Println(fmt.Errorf("error sending to remote webhook url, %w",err))
		}
		defer resp.Body.Close()

		// read the responce
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Println("RESPONCE: "+string(respBody))

		if call.Debug {
			fmt.Println("RESPONCE: "+string(respBody))
		}
	}
}