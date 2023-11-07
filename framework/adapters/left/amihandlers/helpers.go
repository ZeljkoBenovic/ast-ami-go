package amihandlers

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

type direction string

const (
	outbound direction = "outbound"
	inbound  direction = "inbound"
)

func encodeMap(m map[string]string, toFile string) error {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not marshal map: %w", err)
	}

	// write to file if enabled otherwise print to console
	if toFile != "" {
		file, err := os.OpenFile(toFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("could not create %s file: %w", toFile, err)
		}

		defer file.Close()

		if _, err = file.Write(jsonStr); err != nil {
			return fmt.Errorf("could not write %s file: %w", toFile, err)
		}

		_, _ = file.WriteString("\n")
	} else {
		fmt.Printf("%s\n", string(jsonStr))
	}

	return nil
}

func (a *Adapter) sendDataToWebhook(uid string, way direction) {
	if way == inbound {
		if err := a.webhook.SendToWebhook(a.amiEvents.Inbound[CallUID(uid)]); err != nil {
			a.logger.Error("could not send data to the webhook",
				"url", a.config.WebhookURL,
				"method", a.config.WebhookMethod,
				"call_direction", string(inbound),
				"err", err.Error())
		}
	} else {
		if err := a.webhook.SendToWebhook(a.amiEvents.Outbound[CallUID(uid)]); err != nil {
			a.logger.Error("could not send data to the webhook",
				"url", a.config.WebhookURL,
				"method", a.config.WebhookMethod,
				"call_direction", string(outbound),
				"err", err.Error())
		}
	}
}

func convertTimeToUnixTime(timeString string, logger hclog.Logger) int64 {
	timest, err := time.Parse(time.RFC3339Nano, timeString)
	if err != nil {
		logger.Error("Could not parse timestamp", "err", err)
		return 0
	}

	return timest.Unix()
}

func parseAgentName(rawMemberName string, logger hclog.Logger) string {
	match, err := regexp.Match("@", []byte(rawMemberName))
	if err != nil {
		logger.Error("could not use regex on parseAgentName helper func", "err", err.Error())
	}

	if match {
		temp1 := strings.Split(rawMemberName, "@")
		temp2 := strings.Split(temp1[0], "/")

		return temp2[1]
	}

	return rawMemberName
}

func normalizeNumber(number string) string {
	// normalize number only if there is one leading zero
	if string(number[0]) == "0" && string(number[1]) != "0" {
		return "+381" + number[1:]
	}

	return number
}
