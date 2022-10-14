package amihandlers

import (
	"github.com/ZeljkoBenovic/ast-ami-go/framework/ports"
	"github.com/ZeljkoBenovic/ast-ami-go/framework/types/cmd"
	"github.com/hashicorp/go-hclog"
	"github.com/ivahaev/amigo"
)

type Adapter struct {
	amigo   *amigo.Amigo
	config  cmd.Config
	logger  hclog.Logger
	webhook ports.IWebhook

	amiEvents Calls
}

func NewAdapter(config cmd.Config) ports.IAMIHandlers {
	return &Adapter{
		config: config,
		amiEvents: Calls{
			Outbound: map[CallUID]OutboundCall{},
			Inbound:  map[CallUID]InboundCall{},
		},
	}
}

func (a *Adapter) WithAMIGo(amigo *amigo.Amigo) ports.IAMIHandlers {
	a.amigo = amigo

	return a
}

func (a *Adapter) WithWebhook(webhook ports.IWebhook) ports.IAMIHandlers {
	a.webhook = webhook

	return a
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.IAMIHandlers {
	a.logger = logger.Named("ami-handlers")

	return a
}

func (a *Adapter) DebugHandler() {
	if err := a.amigo.RegisterDefaultHandler(func(m map[string]string) {
		if err := encodeMap(m, a.config.AMIDebugFile); err != nil {
			a.logger.Error("Could not encodeMap in DebugHandler", "err", err.Error())
		}
	}); err != nil {
		a.logger.Error("Could not register debug handler", "err", err.Error())
	}
}

func (a *Adapter) RegisterEventHandlers() {
	a.newChannelHandler()
	a.newStateHandler()
	a.queueJoinEvent()
	a.agentConnectEvent()
	a.agentComplete()
	a.queueAbandon()
	a.hangupHandler()
}
