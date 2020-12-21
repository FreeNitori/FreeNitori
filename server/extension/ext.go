package extension

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"io/ioutil"
	"os"
	"strings"
)

var err error
var Commands = make(map[string]string)

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
		if strings.HasSuffix(path.Name(), ".kext") {
			Commands[path.Name()[:len(path.Name())-5]] = "extensions/commands/" + path.Name()
		}
	}
	return nil
}

func RegisterHandlers() error {
	for pattern, path := range Commands {
		multiplexer.ExtensionsCategory.Register(func(context *multiplexer.Context) {
			// TODO: execute extension
		}, pattern, []string{}, "Extension "+path)
	}
	return nil
}
