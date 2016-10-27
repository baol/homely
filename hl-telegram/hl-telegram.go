// homely-telegram
package main

import (
	"encoding/json"
	"flag"
	"log"

	homely "github.com/baol/homely/lib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	bot "github.com/meinside/telegram-bot-go"
)

const (
	telegramVerbose = false
)

// Go on Telegram and grab your token from the @BotFather
func makeTelegramClient(apiToken *string) *bot.Bot {
	client := bot.NewClient(*apiToken)
	client.Verbose = telegramVerbose
	if me := client.GetMe(); !me.Ok {
		panic("Failed to initialize telegram bot")
	}
	return client
}

func telegramChannel(client *bot.Bot, userID *int64) chan mqtt.Message {
	c := make(chan mqtt.Message)
	go func() {
		for {
			msg := <-c
			var message homely.NotificationMessage
			log.Println(msg.Topic(), string(msg.Payload()))
			if err := json.Unmarshal(msg.Payload(), &message); err != nil {
				log.Println("Failed to decode message")
			}
			log.Printf(message.Text)
			options := make(map[string]interface{})
			if sent := client.SendMessage(*userID, &message.Text, options); !sent.Ok {
				log.Printf("Failed to send message: %s\n", *sent.Description)
			}
		}
	}()
	return c
}

func main() {
	log.SetPrefix("hl-telegram: ")
	apiToken := flag.String("telegram-key", "", "Telegram bot key obtained from the @BotFather")
	userID := flag.Int64("user-id", 0, "Used id to be contacted (TODO: document how to get this number!)")
	mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
	flag.Parse()
	required := []string{"telegram-key", "user-id"}
	homely.CheckRequired(required)

	client := makeTelegramClient(apiToken)
	chatChannel := telegramChannel(client, userID)

	queue := mqtt.NewClient(homely.MakeMqttPublishOptions("homely-telegram-"+string(*userID), mqttServer, chatChannel))
	homely.MqttConnectAndSubscribe(queue, map[string]byte{"homely/notification/send": 0})

	homely.MainLoop()
}
