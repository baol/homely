# Homely

## Idea

A collection of glue programs for Domoticz, Frity!Box, and Telegram
using MQTT.

Out of the box, with a few lines of configuration, you can get
notifications for your preferred home automation events on your
telegram account.

E.g. "Baol's back home", "Main door opened", and so on...

But this is also an enabler for more MQTT fun!

## Implementation

### Working

#### hl-telegram

This bot will listen to homely/telegram/send and send the received
message using the configured account and userid.

    {
        "message": "Main door open"
    }

### TODO

#### hl-telegram

The bot will also publish the messages received to homely/telegram/in
using the same format using webhooks.

    {
        "message": "Watch TV"
    }


#### hl-fritzwho

This bot will poll the Fritz!BOX API and check for added/removed
connected devices.

Every time a device gets connected or disconnected it will publish a
message to homely-fritzwho/out/MacAddress/Status, e.g.

    homely-fritzwho/out/PhoneMacAddress/Connected

Useful to know witch devices are active (e.g. phones) for automatic
presence notification.

#### hl-domofilter

Domofilter listens to domoticz/out messages republishes them to
homely/status/${device-id}/${value} when the status changes.

        homely/status/24/On
        homely/status/24/Off

Makes it easier to use *wiring* and *telegram* together.

Analgously devices can be controlled sending messages to
homely/command/24/On

#### hl-wiring

Wiring will listen on multiple MQTT queues and republish the messages
into other queues, after filtering and applying a transformation.

    source-queue dest-queue json

Some form of xpath filtering and json transformation (e.g. JavaScript
will be needed) as well in the rules.

Wiring is stateless, to implement stateful actions we will need
another bot to accumulate state and emit events to be used by
wiring/domofilter/etc.

#### hl-hsm

A state machine for MQTT that follows a flow chart.
