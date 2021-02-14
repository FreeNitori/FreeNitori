// +build linux freebsd darwin

package discord

import (
	"errors"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"io/ioutil"
	"os"
	"plugin"
	"strings"
)

func loadPlugins() error {
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
		symbol, err := pl.Lookup("CommandRoute")
		if err != nil {
			continue
		}
		route, ok := symbol.(*multiplexer.Route)
		if !ok {
			log.Warnf("No Route found in %s.", path.Name())
			continue
		}
		multiplexer.Router.Route(route)
		log.Infof("Loaded plugin %s implementing command %s.", path.Name(), route.Pattern)
	}
	return nil
}
