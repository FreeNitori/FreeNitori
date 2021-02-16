// Package db provides a wrapper around database backend.
package db

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
)

var err error

// Initialize prepares database.
func Initialize() error {
	// Open the database
	err = database.Database.Open(config.Config.System.Database)
	if err != nil {
		return err
	}

	return nil
}
