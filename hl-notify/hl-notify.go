// -*- mode: go; tab-width: 2 -*-
package main

import (
	//"encoding/json"
	"encoding/json"
	"flag"
	//"fmt"
	"log"

	homely "github.com/baol/homely/lib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func notify(notificationChannel <-chan mqtt.Message) {
	for {
		msg := <-notificationChannel
		var message homely.NotificationMessage
		if err := json.Unmarshal(msg.Payload(), &message); err != nil {
			log.Printf("Failed to decode message")
		}
		log.Printf(message.Text)
		homely.LocalNotification(message.Text)
	}
}

func main() {
	log.Println("Starting")
	mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
	flag.Parse()

	notificationChannel := make(chan mqtt.Message)
	queue := mqtt.NewClient(homely.MakeMqttPublishOptions("homely-notification", mqttServer, notificationChannel))
	homely.MqttConnectAndSubscribe(queue, map[string]byte{"homely/notification/send": 0})
	go notify(notificationChannel)

	homely.MainLoop()
}
