package webhook

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

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
	// setting additional headers
	a.setAdditionalHeaders(req)
	// print request headers in debug log
	a.debugHTTPRequestHeaders(req)

	// set new http client with secure or insecure TLS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//nolint
				InsecureSkipVerify: a.config.TLSInsecure,
			},
		},
	}

	// send the request and wait for response
	resp, respErr := client.Do(req)
	if respErr != nil {
		return fmt.Errorf("error sending to remote webhook url, %w", respErr)
	}

	defer resp.Body.Close()

	// print response headers in debug log
	a.debugHTTPRequestHeaders(resp)

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

func (a *Adapter) setAdditionalHeaders(req *http.Request) {
	// if there are no additional header configured no need to do anything
	if a.config.AdditionalHeaders == "" {
		a.logger.Debug("no additional headers defined")

		return
	}

	// get headers from additional headers flag
	headers := strings.Split(a.config.AdditionalHeaders, ",")
	// and add each header to the request
	for _, singleHeader := range headers {
		// first part is header name, second part is header value
		header := strings.Split(singleHeader, ":")
		a.logger.Debug("setting additional header", "header_name", header[0], "header_value", header[1])
		// trim all space and add header to the request
		req.Header.Set(strings.TrimSpace(header[0]), strings.TrimSpace(header[1]))
	}
}

func (a *Adapter) debugHTTPRequestHeaders(req interface{}) {
	switch r := req.(type) {
	case *http.Request:
		reqDump, err := httputil.DumpRequestOut(r, false)
		if err != nil {
			a.logger.Debug("could not dump request headers", "err", err.Error())
		}

		a.logger.Debug("REQUEST HEADERS DUMP", "headers", string(reqDump))
	case *http.Response:
		respDump, err := httputil.DumpResponse(r, false)
		if err != nil {
			a.logger.Debug("could not dump response headers", "err", err.Error())
		}

		a.logger.Debug("RESPONSE HEADERS DUMP", "headers", string(respDump))
	}
}
