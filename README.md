# Homely

## Introduction

Unix style IoT with MQTT.

Instead of developing YAMS (yet another monolithic solution) I'm
trying to develop a bunch of small and independent softwares that
communicate using MQTT that I call *homely*.

This tools are written in [golang](http://golang.org) so they will not
kill your Raspberry PI and should also be reasonably easy to adapt and
extend in case they do not fit your exact needs.

## Installation

You need to install [golang](http://golang.org) first, then you can
use the go tool to install `homely` with all the needed dependency
with the following lines:

    export GOPATH=~/go
    export PATH=$PATH:$GOPATH/bin

    go get -u github.com/baol/homely/...
    go install github.com/baol/homely/...

## Prerequisites

A Unix background :)

A Linux computer (e.g. Raspberry PI) with some kind or IoT devices
attached (433MHz, Z-Wave, ZigBee, or whatever).

An MQTT broker running on one of your machines (the reference
implementation available at https://mosquitto.org/ will do, as should
any packaged version).

In the following examples we will assume that the machine running MQTT
is reachable at the address `mqtt.local`.

Recently services like https://www.cloudmqtt.com/ started to happear
and they offer a free plan if you need to reach your broker from the
public internet, but for our examples installing mosquitto on your
Raspberry PI should be enough.

## Example Deployment Diagram

![Example Deployment Diagram](homely-diagram.png)

## Examples

### Receive a desktop notification when the main door opens

**A Domoticz user wants to receive Desktop notifications every time
the Main door of his apartment opens.**

For this we will need the hl-telegram, hl-notify, hl-wiring and
hl-domofilter modules.

First we need to configure domoticz to forward his messages to MQTT,
so we go to the Hardware section and add a new "MQTT Client Gateway
with LAN Interface" configured to forward all the messages to
`homely.local` in the Flat (default) format.

Unfortunately Domoticz choice of topics does not fit well with homely
so we need to run

    hl-domofilter --mqtt tcp://mqtt.local:1883

To process some interesting Domoticz events and republish them as
homely events. Assuming our main door sensor has id 2 (look in
Domoticz Devices to know the id you are interested in), when the door
will open this will publish a message to `homely/status/2/On`.

We also need to run the notifier

    hl-notify --mqtt tcp://mqtt.local:1883

that will listen `homely/notify/send` for messages in the format
`{"message": "Main door open"}` and send them to your desktop on Mac
OS or Linux.

In order to wire the two together we also need `hl-wiring`

Edit `~/.homely/wiring.toml` and write your rule there using the same
id mentioned above:

    [rule."homely/status/2/On"."homely/notification/send"]
    payload = '{"message": "Main door open"}'

Now launch

    hl-wiring --mqtt tcp://mqtt.local:1883

And enjoy your notifications!

## Other tools

* `hl-telegram` works the same way as `hl-notify`, but you first need
  to register your bot with the @BotFather and know the numeric userid
  you want to send notifications to.
  In order to know your id, send a message to your bot and go to

        https://api.telegram.org/bot<BOT-KEY>/getUpdates

  (replace `<BOT-KEY>` with the key received from the BotFather)

* `hl-flag` works together with an Arduino powered physical
  notificatin devices to raise a flag on certain events, and can be
  controlled sending empty messages to `homely/flag/up` and
  `homely/flag/down`. See
  [Materia Flag](https://github.com/baol/homely/tree/master/hl-flag/materia-flag)
  for details about the device.

  To wire your main door events to the Materia Flag add the following
  rules in `wiring.toml`

        [rule."homely/status/2/Off"."homely/flag/down"]
        [rule."homely/status/2/On"."homely/flag/up"]

Wiring can also be used to automate other actions, like switching off
the lights when you turn on the TV: all messages sent to
`homely/command/<ID>/On` and `homely/command/<ID>/Off` will be translated
in Domoticz commands by `hl-domofilter`.

Happy hacking!
