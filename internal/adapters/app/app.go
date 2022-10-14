package app

import (
	"os"
	"os/signal"
	"syscall"

	frameworkPorts "github.com/ZeljkoBenovic/ast-ami-go/framework/ports"
	internalPorts "github.com/ZeljkoBenovic/ast-ami-go/internal/ports"
	"github.com/hashicorp/go-hclog"
)

type Adapter struct {
	logger      hclog.Logger
	core        internalPorts.ICorePort
	amihandlers frameworkPorts.IAMIHandlers
	webhook     frameworkPorts.IWebhook
	closeCh     chan os.Signal
}

func NewAdapter(
	core internalPorts.ICorePort,
	amihandlers frameworkPorts.IAMIHandlers,
	webhook frameworkPorts.IWebhook,
) internalPorts.IAppPort {
	return &Adapter{
		core:        core,
		amihandlers: amihandlers,
		webhook:     webhook,
		closeCh:     make(chan os.Signal, 1),
	}
}

func (a *Adapter) WithLogger(logger hclog.Logger) internalPorts.IAppPort {
	a.logger = logger.Named("app")

	return a
}

func (a *Adapter) Run() {
	// handle terminate signals
	signal.Notify(a.closeCh, syscall.SIGINT, syscall.SIGTERM)
	// connect to Asterisk AMI
	a.core.ConnectToAsterisk()

	// connect amihandlers and core
	a.amihandlers.WithAMIGo(a.core.AMI())
	a.amihandlers.WithWebhook(a.webhook)

	// if ami-debug flag enabled output all AMI envents
	if a.core.Config().AMIDebug {
		a.amihandlers.DebugHandler()
	}

	// register required event handlers
	a.amihandlers.RegisterEventHandlers()

	// exit on terminate signal
	<-a.closeCh
	a.logger.Info("Received shutdown signal. Program terminated...")
}
