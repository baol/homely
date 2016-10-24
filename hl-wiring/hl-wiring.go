// Homely Wiring
//
// 1. define rules in toml format
// 2. run
//
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/baol/homely/lib"
	"github.com/eclipse/paho.mqtt.golang"
)

type RulePayload struct {
	Payload string
}

type Config struct {
	Rules map[string]map[string]RulePayload `toml:"rule"`
}

func republish(channel chan mqtt.Message, queue mqtt.Client, config Config) {
	for {
		msg := <-channel
		topic := msg.Topic()
		for k, v := range config.Rules[topic] {
			fmt.Println("Sending:", k, v)
			if token := queue.Publish(k, 0, false, v.Payload); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}
	}
}

func main() {
	mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
	flag.Parse()

	var config Config

	if _, err := toml.DecodeFile(os.ExpandEnv("${HOME}/.homely/wiring.toml"), &config); err != nil {
		panic(err)
	}

	fmt.Println(config)

	topics := make(map[string]byte)
	for k := range config.Rules {
		topic := k //strings.Replace(k, "_", "/", -1)
		fmt.Println("Subscribe:", topic)
		topics[topic] = 0
	}

	channel := make(chan mqtt.Message)
	queue := mqtt.NewClient(homely.MakeMqttPublishOptions("hl-wiring", mqttServer, channel))
	go republish(channel, queue, config)
	homely.MqttConnectAndSubscribe(queue, topics)

	homely.MainLoop()
}
