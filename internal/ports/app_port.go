package ports

import "github.com/hashicorp/go-hclog"

type IAppPort interface {
	WithLogger(logger hclog.Logger) IAppPort
	Run()
}
