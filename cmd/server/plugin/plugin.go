// +build linux freebsd darwin

package plugin

import (
	"errors"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database/badger"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"io/ioutil"
	"os"
	"plugin"
	"strings"
)

// Initialize loads all plugins.
func Initialize() error {
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
		symbol, err := pl.Lookup("Setup")
		if err != nil {
			log.Warnf("No setup symbol found in plugin %s.", path.Name())
			continue
		}
		setup, ok := symbol.(func() interface{})
		if !ok {
			log.Warnf("Invalid plugin: setup function found in %s has incorrect function signature.", path.Name())
			continue
		}
		processReturn(setup(), path)
	}
	processReturnPost()
	return nil
}

func processReturn(i interface{}, path os.FileInfo) {
	if db, ok := i.(database.Backend); ok {
		if database.Database != nil {
			log.Warnf("Already loaded database backend %s, skipping plugin %s.", database.Database.DBType(), path.Name())
			return
		}
		database.Database = db
		log.Infof("Loaded plugin %s implementing database backend %s.", path.Name(), db.DBType())
	} else if route, ok := i.(*multiplexer.Route); ok {
		multiplexer.Router.Route(route)
		log.Infof("Loaded plugin %s implementing command %s.", path.Name(), route.Pattern)
	} else {
		log.Infof("Loaded plugin %s.", path.Name())
	}
}

func processReturnPost() {
	if database.Database == nil {
		log.Info("No database backend loaded from plugins, using built-in database.")
		database.Database = &badger.Database
	}
}
