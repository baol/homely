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

// Config TOML main configuration
// [rule."in/topic"."out/topic"]
// payload=...
//
type Config struct {
	Rules map[string]map[string]RulePayload `toml:"rule"`
}

// RulePayload TOML rule payload to be sent
type RulePayload struct {
	Payload string `toml:"payload"`
}

func republish(channel chan mqtt.Message, queue mqtt.Client, config Config) {
	for {
		msg := <-channel
		topic := msg.Topic()
		for k, v := range config.Rules[topic] {
			log.Println("Received:", topic, "Sending:", k, v)
			if token := queue.Publish(k, 0, false, v.Payload); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}
	}
}

func main() {

	log.SetPrefix("hl-wiring: ")

	mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
	flag.Parse()

	var config Config
	configFile := os.ExpandEnv("${HOME}/.homely/wiring.toml")
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		panic(err)
	}

	log.Println("Loaded configuration: ", configFile, config)

	topics := make(map[string]byte)
	for k := range config.Rules {
		topic := k
		log.Println("Subscribe:", topic)
		topics[topic] = 0
	}

	channel := make(chan mqtt.Message)
	queue := mqtt.NewClient(homely.MakeMqttPublishOptions(os.ExpandEnv("hl-wiring-${HOSTNAME}"), mqttServer, channel))
	go republish(channel, queue, config)
	homely.MqttConnectAndSubscribe(queue, topics)

	homely.MainLoop()
}
