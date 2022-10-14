package ports

import (
	"github.com/ZeljkoBenovic/ast-ami-go/internal/types"
	"github.com/hashicorp/go-hclog"
	"github.com/ivahaev/amigo"
)

type ICorePort interface {
	// WithLogger loads logger instance
	WithLogger(logger hclog.Logger) ICorePort
	// WithConfig loads configuration instance
	WithConfig(config types.Config) ICorePort
	// ConnectToAsterisk connects to Asterisk AMI interface
	ConnectToAsterisk()
	// AMI exposes core amigo.Amigo instance
	AMI() *amigo.Amigo
	// Config returns a copy of the core configuration parameters
	Config() types.Config
}
