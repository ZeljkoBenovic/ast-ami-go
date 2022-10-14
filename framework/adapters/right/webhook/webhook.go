package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ZeljkoBenovic/ast-ami-go/framework/ports"
	"github.com/ZeljkoBenovic/ast-ami-go/framework/types/cmd"
	"github.com/hashicorp/go-hclog"
)

type Adapter struct {
	config cmd.Config
	logger hclog.Logger
}

func NewAdapter(config cmd.Config) ports.IWebhook {
	return &Adapter{
		config: config,
	}
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.IWebhook {
	a.logger = logger.Named("webhook")

	return a
}

func (a *Adapter) SendToWebhook(dataToSend interface{}) error {
	// marshal data
	jsonEncoded, err := json.Marshal(dataToSend)
	if err != nil {
		return fmt.Errorf("could not marshal data to json: %w", err)
	}

	a.logger.Debug("Data encoded to json", "data", string(jsonEncoded))

	// create new request endpoint
	req, err := http.NewRequest(a.config.WebhookMethod, a.config.WebhookURL, bytes.NewBuffer(jsonEncoded))
	if err != nil {
		return fmt.Errorf("could not create new request endpoint: %w", err)
	}

	// set json header
	req.Header.Set("Content-Type", "application/json")
	// TODO: set optional authentication header - barer token

	// set new http client
	client := &http.Client{}
	// send the request and wait for response
	resp, respErr := client.Do(req)
	if respErr != nil {
		return fmt.Errorf("error sending to remote webhook url, %w", respErr)
	}

	defer resp.Body.Close()

	// log request to file
	a.logger.Info("Data sent to webhook",
		"url", a.config.WebhookURL,
		"method", a.config.WebhookMethod,
		"data", string(jsonEncoded))

	// read the response
	respBody, _ := ioutil.ReadAll(resp.Body)
	// and log it
	a.logger.Info("Response received from webhook",
		"url", a.config.WebhookURL,
		"method", a.config.WebhookMethod,
		"response", string(respBody))

	return nil
}
