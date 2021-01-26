// +build !linux,!freebsd,!darwin

package database

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database/badger"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
)

func loadDatabaseBackend() error {
	log.Info("Plugins are not supported on this platform, using built-in database.")
	vars.Database = &badger.Database
	return nil
}
