package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jessevdk/go-flags"
	"github.com/slack-go/slack"
	"gopkg.in/yaml.v3"
)

type Opts struct {
	Env string `short:"e" long:"env" description:"Name of the destination environment to which you want to send the message"`
	Vars map[string]string `short:"v" long:"vars" description:"Arbitrary variables can be specified (e.g. --vars=fromIP:xxx.xxx.xxx.xxx)"`
	ConfigFile string `short:"f" long:"configfile" description:"You can specify the path to the configuration file. If this option is specified, an error will occur if the configuration file does not exist"`
}

type VarMap map[string]string

type MsgCtx struct {
	SendAt time.Time
	SendTo Environment
	Msg string
	Vars VarMap
}

type Environment struct {
	Adapter string `yaml:"adapter"`
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
}

type Config struct {
	Version      string                 `yaml:"version"`
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
		} else if fileExists("~/.slcknconf") {
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

func buildSlackContents(msgctx *MsgCtx) slack.Message {
	n := msgctx.SendAt

	headerText := slack.NewTextBlockObject("plain_text",
		":eye-in-speech-bubble: Slckn Notify :eye-in-speech-bubble:", false, false)
	headerBlock := slack.NewHeaderBlock(
		headerText,
		slack.HeaderBlockOptionBlockID("notify-header"),
	)

	contextText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf(":calendar: %s | *from* %s", n.Format("2006-01-02 15:04:05"), msgctx.Vars["name"]),
		false, false)
	contextBlock := slack.NewContextBlock("notify-context", contextText)

	dividerBlock := slack.NewDividerBlock()

	mainContentText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("```%s```", msgctx.Msg), false, false)
	mainContentSectionBlock := slack.NewSectionBlock(mainContentText, nil, nil)

	msg := slack.NewBlockMessage(
		headerBlock,
		contextBlock,
		dividerBlock,
		mainContentSectionBlock,
	)

	return msg
}

func sendSlackMessage(config *Config, msgctx *MsgCtx) {
	api := slack.New(msgctx.SendTo.Token)
	channelID, timestamp, err := api.PostMessage(
		msgctx.SendTo.Channel,
		slack.MsgOptionBlocks(buildSlackContents(msgctx).Msg.Blocks.BlockSet...),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

var opts Opts
var parser = flags.NewParser(&opts, flags.Default)

func init() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}

func main() {
	conf, err := loadConfig()
	fmt.Printf("%s", opts.Env)
	if err != nil {
		log.Fatalf("--- error: %v", err)
		fmt.Println(err)
		return
	}
	fmt.Printf("config: %+v", conf)
	fmt.Printf("%s", conf.Environments["default"].Token)

	msgctx := MsgCtx{
		SendAt: time.Now(),
		SendTo: conf.Environments["default"],
		Msg: "yeah??\n\n\n\nYeah!!!!",
		Vars: VarMap{
			"name": "Barian",
		},
	}

	sendSlackMessage(conf, &msgctx)
}
