// -*- mode: go; tab-width: 2 -*-
// +build !darwin

package homely

import (
	"log"
)

// LocalNotification sends a notifications message to the local desktop
func LocalNotification(notification string) {
	// Unsupported platform
	log.Println("UNSUPPORTED PLATFORM: " + notification)
}
