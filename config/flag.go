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
	configFile := flag.String("configFile","","Configuration file in YAML format")
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
		err = yaml.Unmarshal(confBytes, &yamlConfig)
		if err != nil {
			return config{}, errors.New("could not unmarshal config file")
		}
		
		defer buff.Close()
		
		return yamlConfig, nil
	}

}