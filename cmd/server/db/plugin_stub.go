// +build !linux,!freebsd,!darwin

package db

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database/badger"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
)

func loadDatabaseBackend() error {
	log.Info("Plugins are not supported on this platform, using built-in database.")
	database.Database = &badger.Database
	return nil
}
