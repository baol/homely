// Package homely                          -*- mode: go; tab-width: 2 -*-
//
// Common library for homely daemons
//
package homely

import (
  "encoding/json"
  "flag"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"

  "github.com/eclipse/paho.mqtt.golang"
  bot "github.com/meinside/telegram-bot-go"
)

const (
  // TelegramVerbose controls verboity of the telegram bot
  TelegramVerbose = false
)

// Message mqtt queue communication type
type Message struct {
  Message string `json:"message"`
}

// MakeTelegramClient creates a telegram client given the token
// obtained from the @BotFather
func MakeTelegramClient(apiToken *string) *bot.Bot {
  client := bot.NewClient(*apiToken)
  client.Verbose = TelegramVerbose
  if me := client.GetMe(); !me.Ok {
    panic("Failed to initialize telegram bot")
  }
  return client
}

// TelegramChannel creates a channel and listens on it for new strings
// to be sent to userId
func TelegramChannel(client *bot.Bot, userID *int64) chan string {
  c := make(chan string)
  go func() {
    for {
      message := <-c
      options := make(map[string]interface{})
      if sent := client.SendMessage(*userID, &message, options); !sent.Ok {
        log.Printf("Failed to send message: %s\n", *sent.Description)
      }
    }
  }()
  return c
}

// MakeMqttOptions inizializes the MQTT client options
func MakeMqttOptions(mqttServer *string, chatChannel chan string) *mqtt.ClientOptions {
  opts := mqtt.NewClientOptions().AddBroker(*mqttServer)
  opts.SetClientID("homely-telegram-bot")
  opts.SetProtocolVersion(3)
  opts.SetDefaultPublishHandler(MakeMessageHandler(chatChannel))
  return opts
}

// MqttConnectAndSubscribe connects and subscribe
func MqttConnectAndSubscribe(queue mqtt.Client) {
  if token := queue.Connect(); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    panic("Cannot connect")
  }

  if token := queue.Subscribe("homely-telegram/out", 0, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    panic("Cannot subscribe")
  }
}

// MakeMessageHandler returns an mqtt handler that receive messages,
// decodes them into Messages and sends them into the given queue
func MakeMessageHandler(c chan string) func(client mqtt.Client, msg mqtt.Message) {
  return func(queue mqtt.Client, msg mqtt.Message) {
    var m Message
    log.Printf(string(msg.Payload()))
    if err := json.Unmarshal(msg.Payload(), &m); err != nil {
      log.Printf("Failed to decode message")
    }
    log.Printf(m.Message)
    c <- m.Message
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
