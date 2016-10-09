package main

import (
  "os"
  "os/signal"
  "syscall"
  "fmt"
  "log"
  "flag"
  mqtt "github.com/eclipse/paho.mqtt.golang"
  bot "github.com/meinside/telegram-bot-go"
)

const (
  Verbose = false
)


func telegramChannel(client * bot.Bot, userId * int64)  chan string{
  c := make(chan string)
  go func() {
    for {
      message := <- c
      options := make(map[string]interface{})
      if sent := client.SendMessage(*userId, &message, options); !sent.Ok {
        log.Printf("Failed to send message: %s\n", *sent.Description)
      }
    }
  }()
  return c
}

func makeMessageHandler(c chan string) func (client mqtt.Client, msg mqtt.Message) {
  return func (queue mqtt.Client, msg mqtt.Message) {
    message := fmt.Sprintf("TOPIC: %s\nPAYLOAD: %s\n", msg.Topic(), msg.Payload())
    c <- message
  }
}

func mainLoop() {
  exitSignal := make(chan os.Signal)
  signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
  <-exitSignal
}

func checkRequired(){
  required := []string{"telegram-key", "default-user-id"}
  seen := make(map[string]bool)
  flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
  for _, req := range required {
    if !seen[req] {
      fmt.Fprintf(os.Stderr, "missing required -%s argument\n", req)
      os.Exit(2)
    }
  }
}

func main() {

  apiToken := flag.String("telegram-key", "", "Telegram bot key obtained from the @BotFather")
  userId := flag.Int64("default-user-id", 0, "Used id to be contected by default")
  mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")

  flag.Parse()

  client := bot.NewClient(*apiToken)
  client.Verbose = Verbose

  if me := client.GetMe(); !me.Ok {
    panic("Failed to initialize telegram bot")
  }

  chatChannel := telegramChannel(client, userId)

  opts := mqtt.NewClientOptions().AddBroker(*mqttServer)
  opts.SetClientID("homely-telegram-bot")
  opts.SetProtocolVersion(3)
  opts.SetDefaultPublishHandler(makeMessageHandler(chatChannel))

  queue := mqtt.NewClient(opts)

  if token := queue.Connect(); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    panic("Cannot connect")
  }

  if token := queue.Subscribe("homely-telegram/out", 0, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    panic("Cannot subscribe")
  }

  chatChannel <- "Goodmorning!"

  mainLoop()
}
