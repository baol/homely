// Homely Wiring
//
// 1. define rules in toml format
// 2. run
//
package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/baol/homely/lib"
	"github.com/eclipse/paho.mqtt.golang"
)

type rulePayload struct {
	payload string
}

type config struct {
	rules map[string]map[string]rulePayload `toml:"rule"`
}

func republish(channel chan mqtt.Message, queue mqtt.Client, cfg config) {
	for {
		msg := <-channel
		topic := msg.Topic()
		for k, v := range cfg.rules[topic] {
			log.Println("Received:", topic, "Sending:", k, v)
			if token := queue.Publish(k, 0, false, v.payload); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}
	}
}

func main() {

	log.SetPrefix("hl-wiring: ")

	mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
	flag.Parse()

	var cfg config

	if _, err := toml.DecodeFile(os.ExpandEnv("${HOME}/.homely/wiring.toml"), &cfg); err != nil {
		panic(err)
	}

	log.Println(cfg)

	topics := make(map[string]byte)
	for k := range cfg.rules {
		topic := k
		log.Println("Subscribe:", topic)
		topics[topic] = 0
	}

	channel := make(chan mqtt.Message)
	queue := mqtt.NewClient(homely.MakeMqttPublishOptions(os.ExpandEnv("hl-wiring-${HOSTNAME}"), mqttServer, channel))
	go republish(channel, queue, cfg)
	homely.MqttConnectAndSubscribe(queue, topics)

	homely.MainLoop()
}
