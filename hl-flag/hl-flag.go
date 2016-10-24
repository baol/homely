// homely-flag
//
// Listens on MQTT and interacts with Materia Flag (see materia-flag
// subfolder for the arduino project and software)
//
package main

import (
	"flag"
	"io"

	"github.com/baol/homely/lib"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/jacobsa/go-serial/serial"
)

// Write messages on the serial port for our events
func intercept(c chan mqtt.Message, port io.ReadWriteCloser) {
	for {
		msg := <-c
		switch msg.Topic() {
		case "homely/flag/up":
			port.Write([]byte("110\n")) // flag up at 110°
		case "homely/flag/down":
			port.Write([]byte("0\n")) //   flag down at 0°
		}
	}
}

func main() {
	options := serial.OpenOptions{
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	mqttServer := flag.String("mqtt", "tcp://localhost:1883", "MQTT address")
	flagPort := flag.String("materia", "/dev/ttyUSB0", "USB port for Materia Flag")
	flag.Parse()

	options.PortName = *flagPort

	if port, err := serial.Open(options); err != nil {
		panic(err)
	} else {
		defer port.Close()
		channel := make(chan mqtt.Message)
		queue := mqtt.NewClient(homely.MakeMqttPublishOptions("hl-domoticz", mqttServer, channel))
		homely.MqttConnectAndSubscribe(queue, map[string]byte{"homely/flag/#": 0})

		go intercept(channel, port)

		homely.MainLoop()
	}
}
