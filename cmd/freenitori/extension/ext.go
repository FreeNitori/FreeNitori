package extension

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"io/ioutil"
	"os"
	"strings"
)

var hasExtensions = false

// ExtensionsCategory is a category for commands loaded from extensions.
var ExtensionsCategory = multiplexer.NewCategory("Extensions",
	"Commands loaded in as extensions.")

// Setup sets up extensions.
func Setup() error {
	// Create directories if not exists
	if _, err := os.Stat("extensions"); os.IsNotExist(err) {
		err = os.Mkdir("extensions", 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("extensions/commands/"); os.IsNotExist(err) {
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
				state.Multiplexer.Route(route)
				log.Infof("Loaded JavaScript extension %s.", path.Name())
			}
		}
	}
	if hasExtensions {
		state.Multiplexer.Categories = append(state.Multiplexer.Categories, ExtensionsCategory)
	}
	return nil
}
