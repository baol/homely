// Package homely - common utilities
//
// Common library for homely daemons, mainly converting mqtt topics to go channels
//
package homely

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eclipse/paho.mqtt.golang"
)

// NotificationMessage MQTT message body
type NotificationMessage struct {
	Text string `json:"message"`
}

// MakeMqttOptions inizializes the MQTT client options
func MakeMqttOptions(clientID string, mqttServer *string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker(*mqttServer)
	opts.SetClientID(clientID)
	opts.SetProtocolVersion(3)
	return opts
}

// MakeMqttPublishOptions inizializes the MQTT client options, and sets the publish handler
func MakeMqttPublishOptions(clientID string, mqttServer *string, channel chan mqtt.Message) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker(*mqttServer)
	opts.SetClientID(clientID)
	opts.SetProtocolVersion(3)
	opts.SetDefaultPublishHandler(MakeMessageHandler(channel))
	return opts
}

// MqttConnectAndSubscribe connects and subscribe
func MqttConnectAndSubscribe(queue mqtt.Client, topic map[string]byte) {
	if token := queue.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic("Cannot connect")
	}

	if token := queue.SubscribeMultiple(topic, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic("Cannot subscribe")
	}
}

// MakeMessageHandler returns an mqtt handler that receive messages,
// and sends them to the channel
func MakeMessageHandler(c chan mqtt.Message) func(client mqtt.Client, msg mqtt.Message) {
	return func(queue mqtt.Client, msg mqtt.Message) {
		c <- msg
	}
}

// CheckRequired checks for required flags on the command line
func CheckRequired(required []string) {
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s argument\n", req)
			os.Exit(2)
		}
	}
}

// MainLoop waits until a SIGTERM (e.g. Ctrl-C) is received
func MainLoop() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)
	<-exitSignal
}
