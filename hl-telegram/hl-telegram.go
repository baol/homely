// homely-telegram
package main

import (
  "flag"

  homely "github.com/baol/homely/lib"
  mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
  apiToken := flag.String("telegram-key", "", "Telegram bot key obtained from the @BotFather")
  userID := flag.Int64("default-user-id", 0, "Used id to be contected by default")
  mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
  flag.Parse()
  required := []string{"telegram-key", "default-user-id"}
  homely.CheckRequired(required)

  client := homely.MakeTelegramClient(apiToken)
  chatChannel := homely.TelegramChannel(client, userID)

  queue := mqtt.NewClient(homely.MakeMqttOptions("homely-telegram", mqttServer, chatChannel))
  homely.MqttConnectAndSubscribe(queue, map[string]byte{"homely/telegram/send": 0})

  homely.MainLoop()
}
