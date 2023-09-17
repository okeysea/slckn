package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"github.com/cockroachdb/errors"
)

var data = `
version: 1
environments:
  default:
    adapter: slack-app-token
    token: xoxb-4820273758467-5909894820756-bIHR923k63Ng49KSLIrQgBeZ
    channel: '#general'
  testchannel:
    adapter: slack-app-token
    token: xoxb-4820273758467-5909894820756-bIHR923k63Ng49KSLIrQgBeZ
    channel: '#general'
`

type Environment struct {
	Adapter string `yaml:"adapter"`
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
}

type Config struct {
	Version      string `yaml:"version"`
	Environments map[string]Environment `yaml:"environments"`
}

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}

func loadFile(filepath string) ([]byte, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "Can not read file")
	}

	return bytes, nil
}

func loadConfigFromByte(confbyte []byte) (*Config, error) {
	config := Config{}
	err := yaml.Unmarshal(confbyte, &config)
	if err != nil {
		return nil, errors.Wrap(err, "Can not load config")
	}

	return &config, nil
}

func loadConfig() (*Config, error) {
	path := os.Getenv("SLCKN_CONFIG_PATH")
	if path == "" {
		if fileExists("./.slcknconf") {
			path = "./.slcknconf"
		}else if fileExists("~/.slcknconf") {
			path = "~/.slcknconf"
		}
	}

	if !fileExists(path) {
		return nil, errors.New("Config file not found") 
	}

	bytes, err := loadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Load config error")
	}

	conf, err := loadConfigFromByte(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Parse config error")
	}

	return conf, nil
}

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("--- error: %v", err)
		fmt.Println(err)
		return
	}
	fmt.Printf("config: %+v", conf)

	fmt.Printf("%s", conf.Environments["default"].Token)
}
