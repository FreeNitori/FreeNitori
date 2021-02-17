package extension

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"io/ioutil"
	"os"
	"strings"
)

var err error
var hasExtensions = false

// ExtensionsCategory is a category for commands loaded from extensions.
var ExtensionsCategory = multiplexer.NewCategory("Extensions",
	"Commands loaded in as extensions.")

// FindExtensions finds and registers all extensions.
func FindExtensions() error {
	// Create directories if not exists
	_, err = os.Stat("extensions")
	if os.IsNotExist(err) {
		err = os.Mkdir("extensions", 0755)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat("extensions/commands/")
	if os.IsNotExist(err) {
		err = os.Mkdir("extensions/commands/", 0755)
		if err != nil {
			return err
		}
	}

	commandExtensionPaths, err := ioutil.ReadDir("extensions/commands/")
	if err != nil {
		return err
	}
	for _, path := range commandExtensionPaths {
		if strings.HasSuffix(path.Name(), ".njs") {
			route := registerJS("extensions/commands/", path.Name())
			if route != nil {
				hasExtensions = true
				multiplexer.Router.Route(route)
				log.Infof("Loaded JavaScript extension %s.", path.Name())
			}
		}
	}
	if hasExtensions {
		multiplexer.Categories = append(multiplexer.Categories, ExtensionsCategory)
	}
	return nil
}
