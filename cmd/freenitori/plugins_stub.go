// +build !linux,!freebsd,!darwin

package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database/badger"
	log "git.randomchars.net/FreeNitori/Log"
)

func setupPlugins() error {
	log.Info("Plugins are not supported on this platform.")
	database.Database = &badger.Database
	return nil
}
