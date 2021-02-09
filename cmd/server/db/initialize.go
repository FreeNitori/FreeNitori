// Package db provides a wrapper around database backend.
package db

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
)

var err error

// Initialize prepares database.
func Initialize() error {
	// Load database backend
	err = loadDatabaseBackend()
	if err != nil {
		return err
	}

	// Open the database
	err = database.Database.Open(config.Config.System.Database)
	if err != nil {
		return err
	}

	return nil
}