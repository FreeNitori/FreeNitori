// +build linux freebsd darwin

package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database/badger"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"os"
	"plugin"
	"runtime"
	"strings"
)

var plugins []nitoriPlugin

type nitoriPlugin struct {
	plugin   *plugin.Plugin
	Name     string             `json:"name"`
	Valid    bool               `json:"valid"`
	Route    *multiplexer.Route `json:"route"`
	Database database.Backend   `json:"database"`
}

// Lookup searches for a symbol named symName in plugin p.
// A symbol is any exported variable or function.
// It reports an error if the symbol is not found.
// It is safe for concurrent use by multiple goroutines.
func (p *nitoriPlugin) Lookup(symName string) (plugin.Symbol, error) {
	return p.plugin.Lookup(symName)
}

func setupPlugins() error {
	var (
		err error
		ok  bool
	)

	// Read plugin directory, create if necessary
	var dir []os.DirEntry
	if dir, err = os.ReadDir("plugins"); err != nil {
		if os.IsNotExist(err) {
			log.Warn("Creating plugin directory since it does not exist.")
			if err = os.Mkdir("plugins", 0700); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	for _, entry := range dir {
		// Skip entry if no proper suffix
		switch runtime.GOOS {
		case "darwin":
			if !strings.HasSuffix(entry.Name(), ".dylib") {
				continue
			}
		default:
			if !strings.HasSuffix(entry.Name(), ".so") {
				continue
			}
		}

		// Skip if entry is a directory
		if entry.IsDir() {
			continue
		}

		// Load plugin
		var instance *plugin.Plugin
		if instance, err = plugin.Open("plugins/" + entry.Name()); err != nil {
			log.Warnf("Error loading plugin %s, %s", entry.Name(), err)
			continue
		}
		p := nitoriPlugin{
			plugin: instance,
			Name:   entry.Name(),
		}

		// Read setup symbol
		var setup func() interface{}
		var setupSym interface{}
		if setupSym, err = instance.Lookup("Setup"); err != nil {
			log.Warnf("Error looking up setup symbol in plugin %s.", entry.Name())
			// Prematurely append instance to slice
			plugins = append(plugins, p)
			continue
		}
		if setup, ok = setupSym.(func() interface{}); !ok {
			log.Warnf("Invalid setup in plugin %s.", entry.Name())
			p.Valid = false
			continue
		} else {
			p.Valid = true
		}

		// Receive plugin setup return
		var (
			db    database.Backend
			route *multiplexer.Route
		)
		payload := setup()
		// Database
		if db, ok = payload.(database.Backend); ok {
			p.Database = db
			if database.Database != nil {
				p.Valid = false
				log.Warnf("Database backend %s already loaded, skipping plugin %s.",
					database.Database.DBType(), entry.Name())
				continue
			}
			database.Database = db
			log.Infof("Plugin %s loaded implementing database backend %s.", entry.Name(), db.DBType())
			continue
		}
		// Route
		if route, ok = payload.(*multiplexer.Route); ok {
			state.Multiplexer.Route(route)
			p.Route = route
			log.Infof("Plugin %s loaded implementing command route %s.", entry.Name(), route.Pattern)
			continue
		}
		// Error
		if err, ok = payload.(error); ok {
			p.Valid = false
			log.Errorf("Plugin %s error while setting up, %s", entry.Name(), err.Error())
			continue
		}
		// Catch-all
		if payload == nil {
			log.Infof("Plugin %s loaded.", entry.Name())
		} else {
			p.Valid = false
			log.Warnf("Plugin %s unexpected return in setup.", entry.Name())
		}

		// Append instance to slice
		plugins = append(plugins, p)
	}

	// Default database if required
	if database.Database == nil {
		log.Infof("Plugins did not provide any database, falling back to default.")
		database.Database = &badger.Database
	}

	return nil
}
