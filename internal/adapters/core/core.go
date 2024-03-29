package core

import (
	"github.com/ZeljkoBenovic/ast-ami-go/internal/ports"
	"github.com/ZeljkoBenovic/ast-ami-go/internal/types"
	"github.com/hashicorp/go-hclog"
	"github.com/ivahaev/amigo"
)

type Adapter struct {
	logger hclog.Logger
	ami    *amigo.Amigo
	config types.Config
}

func NewAdapter() ports.ICorePort {
	return &Adapter{}
}

func (a *Adapter) WithLogger(logger hclog.Logger) ports.ICorePort {
	a.logger = logger.Named("core")

	return a
}

func (a *Adapter) WithConfig(conf types.Config) ports.ICorePort {
	a.config = conf

	return a
}

func (a *Adapter) AMI() *amigo.Amigo {
	return a.ami
}

func (a *Adapter) Config() types.Config {
	return a.config
}

func (a *Adapter) ConnectToAsterisk() {
	a.ami = amigo.New(&amigo.Settings{
		Username: a.config.Username,
		Password: a.config.Password,
		Host:     a.config.Host,
		Port:     a.config.Port,
	})

	a.ami.Connect()
	a.ami.On("error", func(message string) {
		a.logger.Error("Error received from the Asterisk AMI", "err", message)
	})
	a.ami.On("connect", func(message string) {
		a.logger.Info("Connected to Asterisk AMI",
			"msg", message,
			"host", a.config.Host,
			"port", a.config.Port,
		)
	})
}
