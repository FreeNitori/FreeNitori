// +build linux freebsd darwin

package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"os"
)

func init() {
	// Check for existence of plugin directory
	_, err := os.Stat("plugins")
	if os.IsNotExist(err) {
		err = os.Mkdir("plugins", 0755)
		if err != nil {
			log.Fatalf("Unable to create plugin directory, %s", err)
			os.Exit(1)
		}
	}
}
