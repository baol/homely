// +build !darwin,!linux

package homely

import (
	"log"
)

// LocalNotification sends a notifications message to the local desktop
func LocalNotification(notification string) {
	// Unsupported platform
	log.Println("UNSUPPORTED PLATFORM: " + notification)
}
