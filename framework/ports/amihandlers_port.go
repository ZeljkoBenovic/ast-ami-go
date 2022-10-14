package ports

import (
	"github.com/hashicorp/go-hclog"
	"github.com/ivahaev/amigo"
)

type IAMIHandlers interface {
	// WithLogger loads new logger instance
	WithLogger(logger hclog.Logger) IAMIHandlers
	// WithAMIGo loads existing amigo.Amigo instance that can run methods on
	WithAMIGo(amigo *amigo.Amigo) IAMIHandlers
	// WithWebhook loads new webhook instance
	WithWebhook(webhook IWebhook) IAMIHandlers
	// DebugHandler will output all events received from Asterisk AMI interface.
	// Used for debugging only.
	DebugHandler()
	// RegisterEventHandlers arms all necessary event handlers.
	RegisterEventHandlers()
}
