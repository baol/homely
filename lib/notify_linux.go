// +build linux

package homely

import (
	"log"

	"github.com/godbus/dbus"
)

// LocalNotification sends a notifications message to the local desktop
func LocalNotification(notification string) {
	log.Println(notification)

	if conn, err := dbus.SessionBus(); err != nil {
		log.Println(err)
	} else {
		obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
		call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
			"", "Homely", notification, []string{},
			map[string]dbus.Variant{}, int32(5000))
		if call.Err != nil {
			log.Println(call.Err)
		}
	}
}
