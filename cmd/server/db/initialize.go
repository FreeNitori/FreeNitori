// Package db provides a wrapper around database backend.
package db

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
)

// Initialize prepares database.
func Initialize() error {
	// Open the database
	if err := database.Database.Open(config.Config.System.Database); err != nil {
		return err
	}
	return nil
}
