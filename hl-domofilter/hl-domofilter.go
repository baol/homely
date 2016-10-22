// homely-domofilter
//
// Translates Light/Switch events in two directions:
//
//  * Listens on domoticz/out and republishes them to homely/status/<id>/{On,Off}
//  * Listens on homely/command/<idx>/{On,Off} and sends the appropriate command to domoticz/in
//
// Used to simplify the writing of hl-wiring rules.
//
package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "strings"

  "github.com/baol/homely/lib"
  "github.com/eclipse/paho.mqtt.golang"
)

// republish only works for on/off switches at the moment
func republish(c chan mqtt.Message, queue mqtt.Client) {
  for {
    msg := <-c
    switch msg.Topic() {
    // domoticz to homely
    case "domoticz/out":
      var payload map[string]interface{}
      json.Unmarshal(msg.Payload(), &payload)
      if payload["dtype"] == "Light/Switch" {
        var onoff string
        if (payload["nvalue"].(float64)) == 0 {
          onoff = "Off"
        } else {
          onoff = "On"
        }
        topic := fmt.Sprintf("homely/status/%d/%s", int(payload["idx"].(float64)), onoff)
        fmt.Println(topic)
        if token := queue.Publish(topic, 0, false, msg.Payload()); token.Wait() && token.Error() != nil {
          panic(token.Error())
        }
      }
      // homely to domoticz
    default:
      fmt.Println(msg.Topic())
      tokens := strings.Split(msg.Topic(), "/")
      payload := fmt.Sprintf("{\"command\": \"switchlight\", \"idx\": %s, \"switchcmd\": \"%s\"}", tokens[2], tokens[3])
      fmt.Println("domoticz/in")
      fmt.Println(payload)
      if token := queue.Publish("domoticz/in", 0, false, payload); token.Wait() && token.Error() != nil {
        panic(token.Error())
      }
    }
  }
}

func main() {
  mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
  flag.Parse()

  channel := make(chan mqtt.Message)
  queue := mqtt.NewClient(homely.MakeMqttPublishOptions("hl-domoticz", mqttServer, channel))
  go republish(channel, queue)
  homely.MqttConnectAndSubscribe(queue, map[string]byte{"domoticz/out": 0, "homely/command/#": 0})

  homely.MainLoop()
}
