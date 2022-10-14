package ports

import (
	"github.com/hashicorp/go-hclog"
)

type IWebhook interface {
	SendToWebhook(dataToSend interface{}) error
	WithLogger(logger hclog.Logger) IWebhook
}
