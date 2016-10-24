// -*- mode: go; tab-width: 2 -*-
// +build darwin

package homely

import (
	"log"

	"github.com/deckarep/gosx-notifier"
)

// LocalNotification sends a notifications message to the local desktop
func LocalNotification(notification string) {
	log.Println(notification)
	note := gosxnotifier.NewNotification(notification)
	note.Title = "Homely"
	err := note.Push()
	if err != nil {
		log.Println("Uh oh! Error with Notify")
	}
}
