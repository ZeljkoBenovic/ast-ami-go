package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrNotYamlFile = errors.New("config file must be in YAML format")
)

type Config struct {
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	Host                string `yaml:"host"`
	Port                string `yaml:"port"`
	InboundContext      string `yaml:"inbound_context"`
	OutboundContext     string `yaml:"outbound_context"`
	LogFileLocation     string `yaml:"log_file_location"`
	WebhookURL          string `yaml:"webhook_url"`
	WebhookMethod       string `yaml:"webhook_method"`
	AdditionalHeaders   string `yaml:"additional_headers"`
	LogLevel            string `yaml:"log-level"`
	AMIDebug            bool   `yaml:"ami_debug"`
	AMIDebugFile        string `yaml:"ami_debug_file"`
	MonitorPublicFolder string `yaml:"monitor_public_folder"`

	ExportDefaultConfig bool   `yaml:"-"`
	ConfigFileLocation  string `yaml:"-"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ProcessConfig() error {
	flag.StringVar(&c.Username, "username", "admin", "Username for authentication to the AMI interface")
	flag.StringVar(&c.Password, "password", "amp111", "Password for authentication to the AMI interface")
	flag.StringVar(&c.Host, "host", "127.0.0.1", "AMI server ip address")
	flag.StringVar(&c.Port, "port", "5038", "AMI server port number")

	flag.StringVar(&c.ConfigFileLocation, "config", "", "Configuration file in YAML format")
	flag.StringVar(&c.LogFileLocation, "log-file", "", "Location of log file for call events")

	flag.StringVar(&c.InboundContext, "inbound-context", "from-trunk", "Context for all inbound calls")
	flag.StringVar(&c.OutboundContext, "outbound-context", "from-internal", "Context for all outbound calls")

	flag.StringVar(&c.WebhookURL, "webhook-url", "", "The webhook URL endpoint to send the events to")
	flag.StringVar(&c.WebhookMethod, "webhook-method", "POST", "The REST method used to send the event to webhook")
	flag.StringVar(&c.AdditionalHeaders, "additional-headers", "",
		"Additional headers that will be sent to webhook in format <header_name>:<header_value>, separated by comma (,)")

	flag.BoolVar(&c.AMIDebug, "ami-debug", false, "Set this flag to catch all AMI events")
	flag.StringVar(&c.AMIDebugFile, "ami-debug-file", "", "File to write all AMI events to")

	flag.BoolVar(&c.ExportDefaultConfig, "export", false, "Set this flag to export the default config file")
	flag.StringVar(&c.LogLevel, "log-level", "info", "Turn on the debug mode and output everything to console")

	flag.StringVar(&c.MonitorPublicFolder, "monitor-public", "",
		"The http publicly available folder that the recordings can be downloaded from")

	flag.Parse()

	if err := c.setConfigFromConfigFile(); err != nil {
		return fmt.Errorf("could not load config from file: %w", err)
	}

	return nil
}

func (c *Config) setConfigFromConfigFile() error {
	// do not process if config file is not defined
	if c.ConfigFileLocation == "" {
		return nil
	}

	extension := filepath.Ext(c.ConfigFileLocation)
	if extension != ".yaml" && extension != ".yml" {
		return ErrNotYamlFile
	}

	buff, err := os.Open(c.ConfigFileLocation)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}

	defer buff.Close()

	confBytes, err := ioutil.ReadAll(buff)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	if yamlErr := yaml.Unmarshal(confBytes, c); yamlErr != nil {
		return fmt.Errorf("could not unmarshal config file: %w", err)
	}

	return nil
}
