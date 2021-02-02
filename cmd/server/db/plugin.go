// +build linux freebsd darwin

package db

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database/badger"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"io/ioutil"
	"os"
	"plugin"
	"strings"
)

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
		db, ok := symbol.(database.Backend)
		if !ok {
			log.Warnf("No database backend found in %s.", path.Name())
			continue
		}
		if database.Database != nil {
			log.Warnf("Already loaded database backend %s, skipping plugin %s.", database.Database.DBType(), path.Name())
			continue
		}
		database.Database = db
		log.Infof("Loaded plugin %s implementing database backend %s.", path.Name(), db.DBType())
	}
	if database.Database == nil {
		log.Info("No database backend loaded from plugins, using built-in database.")
		database.Database = &badger.Database
	}
	return nil
}
