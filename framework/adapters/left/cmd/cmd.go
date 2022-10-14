package cmd

import (
	"os"

	"github.com/ZeljkoBenovic/ast-ami-go/framework/ports"
	"github.com/ZeljkoBenovic/ast-ami-go/framework/types/cmd"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/yaml.v3"
)

type Adapter struct {
	logger hclog.Logger
	config *cmd.Config
}

func NewAdapter() ports.ICmdPort {
	return &Adapter{
		config: cmd.NewConfig(),
		logger: hclog.New(&hclog.LoggerOptions{
			Name: "ast-ami-go",
		}),
	}
}

func (a *Adapter) WithLogger() ports.ICmdPort {
	a.logger = a.logger.Named("cmd")

	return a
}

func (a *Adapter) GetConfig() *cmd.Config {
	// process configuration options
	if err := a.config.ProcessConfig(); err != nil {
		a.logger.Error("Could not load configuration", "err", err.Error())
	}

	// set logger options
	a.setLoggerOptions()

	a.logger.Debug("Configuration load success", "config", a.config)

	// export default config file
	if a.config.ExportDefaultConfig {
		a.exportConfig()
		a.logger.Info("Default configuration file exported to ./config.yaml")
		os.Exit(0)
	}

	return a.config
}

func (a *Adapter) exportConfig() {
	yamlBuff, err := yaml.Marshal(a.config)
	if err != nil {
		a.logger.Error("Could not marshal config to buffer", "err", err.Error())
	}

	if writeErr := os.WriteFile("config.yaml", yamlBuff, os.FileMode(0644)); writeErr != nil {
		a.logger.Error("Could not write config file", "err", writeErr.Error())
	}
}

func (a *Adapter) setLoggerOptions() {
	a.logger.SetLevel(hclog.LevelFromString(a.config.LogLevel))

	// set log file if set
	if a.config.LogFileLocation != "" {
		logFile, err := os.Create(a.config.LogFileLocation)
		if err != nil {
			a.logger.Error("Could not write log file", "err", err.Error())
			os.Exit(1)
		}

		a.logger = hclog.New(&hclog.LoggerOptions{
			Name:   "ast-ami-go",
			Output: logFile,
			Level:  hclog.LevelFromString(a.config.LogLevel),
		})

		// remain in the same log context
		a.logger = a.logger.Named("cmd")
	}
}
