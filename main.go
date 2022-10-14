package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ZeljkoBenovic/ast-ami-go/framework/adapters/left/amihandlers"
	"github.com/ZeljkoBenovic/ast-ami-go/framework/adapters/left/cmd"
	"github.com/ZeljkoBenovic/ast-ami-go/framework/adapters/right/webhook"
	"github.com/ZeljkoBenovic/ast-ami-go/internal/adapters/app"
	"github.com/ZeljkoBenovic/ast-ami-go/internal/adapters/core"
	internalTypes "github.com/ZeljkoBenovic/ast-ami-go/internal/types"
	"github.com/hashicorp/go-hclog"
)

func main() {
	// create config
	conf := cmd.NewAdapter().WithLogger().GetConfig()

	// setup logger
	logger, err := newLogger(conf.LogLevel, conf.LogFileLocation)
	if err != nil {
		log.Fatalln("could not create logger instance: ", err.Error())
	}

	// create new main AMI application
	amiApp := app.NewAdapter(
		core.NewAdapter().
			WithLogger(logger).
			WithConfig(internalTypes.Config{
				Username:        conf.Username,
				Password:        conf.Password,
				Host:            conf.Host,
				Port:            conf.Port,
				InboundContext:  conf.InboundContext,
				OutboundContext: conf.OutboundContext,
				LogFileLocation: conf.LogFileLocation,
				WebhookURL:      conf.WebhookURL,
				WebhookMethod:   conf.WebhookMethod,
				LogLevel:        conf.LogLevel,
				AMIDebug:        conf.AMIDebug,
			}),
		amihandlers.NewAdapter(*conf).WithLogger(logger),
		webhook.NewAdapter(*conf).WithLogger(logger),
	).WithLogger(logger)

	// run the application
	amiApp.Run()
}

func newLogger(logLevel, logFile string) (hclog.Logger, error) {
	var (
		logWriter io.Writer = nil
		err       error
	)

	if logFile != "" {
		logWriter, err = os.Create(logFile)
		if err != nil {
			return nil, fmt.Errorf("could not create log file %s err: %w", logFile, err)
		}
	}

	return hclog.New(&hclog.LoggerOptions{
		Name:   "ast-ami-go",
		Output: logWriter,
		Level:  hclog.LevelFromString(logLevel),
	}), nil
}
