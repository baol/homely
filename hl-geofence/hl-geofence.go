// Homely Geofence
//
// Checks weather a Owntracks device is inside a certain radius of the home location
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"strings"

	homely "github.com/baol/homely/lib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type point struct {
	lat float64
	lng float64
}

const (
	earthRadius = 6371
)

// Aversine distance
// http://www.movable-type.co.uk/scripts/latlong.html
func greatCircleDistance(p1, p2 *point) float64 {
	dLat := (p2.lat - p1.lat) * (math.Pi / 180.0)
	dLon := (p2.lng - p1.lng) * (math.Pi / 180.0)

	lat1 := p1.lat * (math.Pi / 180.0)
	lat2 := p2.lat * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func geofence(home point, queue mqtt.Client, notificationChannel <-chan mqtt.Message) {
	fenceCache := make(map[string]int)
	batteryCache := make(map[string]int)
	for {
		msg := <-notificationChannel
		device := strings.Split(msg.Topic(), "/")
		var raw interface{}
		if err := json.Unmarshal(msg.Payload(), &raw); err != nil {
			log.Printf("Failed to decode message")
		}
		message := raw.(map[string]interface{})
		switch message["_type"].(string) {
		case "location":
			loc := point{message["lat"].(float64), message["lon"].(float64)}
			acc := message["acc"].(float64)
			batt := message["batt"].(float64)
			dist := greatCircleDistance(&home, &loc)
			topic := fmt.Sprintf("homely/device/%s/%s", device[1], device[2])
			payload := ""
			if dist > acc+50.0 && fenceCache[topic] >= 0 { // outside
				fenceCache[topic] = -1
				if token := queue.Publish(topic+"/away", 0, false, payload); token.Wait() && token.Error() != nil {
					panic(token.Error())
				}

			} else if dist < acc+20.0 && fenceCache[topic] <= 0 { // inside
				fenceCache[topic] = +1
				if token := queue.Publish(topic+"/home", 0, false, payload); token.Wait() && token.Error() != nil {
					panic(token.Error())
				}
			}
			payload = fmt.Sprintf("{\"battery\": %f }", batt)
			if batt < 30.0 && batteryCache[topic] >= 0 {
				batteryCache[topic] = -1
				if token := queue.Publish(topic+"/battery/low", 0, false, payload); token.Wait() && token.Error() != nil {
					panic(token.Error())
				}
			} else if batteryCache[topic] <= 0 {
				batteryCache[topic] = +1
				if token := queue.Publish(topic+"/battery/good", 0, false, payload); token.Wait() && token.Error() != nil {
					panic(token.Error())
				}
			}
		}
	}
}

func main() {
	mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
	homeLat := flag.Float64("lat", 0, "Latitude of your place")
	homeLon := flag.Float64("lon", 0, "Longitude of your place")
	flag.Parse()
	required := []string{"lat", "lon"}
	homely.CheckRequired(required)

	notificationChannel := make(chan mqtt.Message)
	queue := mqtt.NewClient(homely.MakeMqttPublishOptions("homely-owntracks", mqttServer, notificationChannel))
	homely.MqttConnectAndSubscribe(queue, map[string]byte{"owntracks/#": 0})
	go geofence(point{*homeLat, *homeLon}, queue, notificationChannel)

	homely.MainLoop()
}
