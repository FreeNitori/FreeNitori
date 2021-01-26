// Wrapper around database backend driver.
package database

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
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
	err = vars.Database.Open(config.Config.System.Database)
	if err != nil {
		return err
	}

	return nil
}
