// Wrapper around database backend driver.
package database

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	"io/ioutil"
	"os"
	"plugin"
	"strings"
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

// loadDatabaseBackend loads database backend from plugins.
func loadDatabaseBackend() error {
	stat, err := os.Stat("plugins")
	if os.IsNotExist(err) {
		return errors.New("plugins directory does not exist")
	}
	if !stat.IsDir() {
		return errors.New("plugins path is not a directory")
	}
	pluginPaths, err := ioutil.ReadDir("plugins/")
	if err != nil {
		return errors.New("plugins directory unreadable")
	}
	for _, path := range pluginPaths {
		if !strings.HasSuffix(path.Name(), ".so") {
			continue
		}
		pl, err := plugin.Open("plugins/" + path.Name())
		if err != nil {
			log.Warnf("Error while loading plugin %s, %s", path.Name(), err)
			continue
		}
		symbol, err := pl.Lookup("Database")
		if err != nil {
			continue
		}
		db, ok := symbol.(vars.Backend)
		if !ok {
			log.Warnf("No database backend found in %s.", path.Name())
			continue
		}
		if vars.Database != nil {
			log.Warnf("Already loaded database backend %s, skipping plugin %s.", vars.Database.DBType(), path.Name())
			continue
		}
		vars.Database = db
		log.Infof("Loaded plugin %s implementing database backend %s.", path.Name(), db.DBType())
	}
	if vars.Database == nil {
		return errors.New("no database backend loaded from plugins")
	}
	return nil
}
