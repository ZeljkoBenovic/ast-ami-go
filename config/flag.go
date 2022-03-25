package config

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	InboundContext string `yaml:"inbound_context"`
	OutboundContext string `yaml:"outbound_context"`
	LogFileLocation string `yaml:"log_file_location"`
	WebhookURL string `yaml:"webhook_url"`
	WebhookMethod string `yaml:"webhook_method"`
	Debug bool `yaml:"debug"`
}

func GetConfig() config {

	config, err := processFlags(os.Args[1:])

	if err != nil {
		fmt.Println("Could not get config. Error: ", err.Error())
		os.Exit(1)
	}	
	
	return config 

}

func processFlags(args []string) (config, error) {

	username := flag.String("username", "admin", "Username for authentication to the AMI interface")
	password := flag.String("password", "amp111", "Password for authentication to the AMI interface")
	host := flag.String("host", "127.0.0.1", "AMI server ip address")
	port := flag.String("port", "5038", "AMI server port number")
	configFile := flag.String("config","","Configuration file in YAML format")

	inboundContext := flag.String("inbound-context", "from-trunk", "Context for all inbound calls")
	outboundContext := flag.String("outbound-context", "from-internal", "Context for all outbound calls")

	logFileLocation := flag.String("log-file","/tmp/call_event.log","Location of log file for call events")

	webhookUrl := flag.String("webhook-url", "","The webhook URL endpoint to send the events to")
	webhookMethod := flag.String("webhook-method","POST","The REST method used to send the event to webhook")

	debug := flag.Bool("debug", false, "Turn on the debug mode and output everything to console")
	flag.Parse()

	if len(args) == 0 {
		fmt.Println("Using default username, password and host. (admin, admin, 127.0.0.1)")
		fmt.Println("To see available options use -h or --help flag.")
	}

	if *configFile == "" {
		return config{
			Username: *username,
			Password: *password, 
			Host: *host,
			Port: *port,
			InboundContext: *inboundContext,
			OutboundContext: *outboundContext,
			LogFileLocation: *logFileLocation,
			WebhookURL: *webhookUrl,
			WebhookMethod: *webhookMethod,
			Debug: *debug,
			}, nil
	} else {
		extension := filepath.Ext(*configFile)
		if ( extension != ".yaml" || extension == ".yml" ) {
			return config{}, errors.New("config file must be in YAML format")
		}

		buff, err := os.Open(*configFile)
		if err != nil {
			return config{}, errors.New("could not open config file")
		}
		
		confBytes, err := ioutil.ReadAll(buff)
		if err != nil {
			return config{}, errors.New("could not read config file")
		}
		
		yamlConfig := config{}
		if err := yaml.Unmarshal(confBytes, &yamlConfig); err != nil {
			return config{}, errors.New("could not unmarshal config file")
		}
		
		defer buff.Close()
		
		return yamlConfig, nil
	}

}