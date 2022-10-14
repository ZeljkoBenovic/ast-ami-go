package ports

import (
	"github.com/ZeljkoBenovic/ast-ami-go/framework/types/cmd"
)

type ICmdPort interface {
	WithLogger() ICmdPort
	GetConfig() *cmd.Config
}
