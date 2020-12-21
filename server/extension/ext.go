package extension

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"io/ioutil"
	"os"
	"strings"
)

var err error
var Commands = make(map[string]string)
var ExtensionsCategory = multiplexer.NewCategory("Extensions",
	"Commands loaded in as extensions.")

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
	if len(Commands) > 0 {
		multiplexer.Categories = append(multiplexer.Categories, ExtensionsCategory)
	}
	return nil
}

func RegisterHandlers() error {
	for pattern, path := range Commands {
		multiplexer.Router.Route(&multiplexer.Route{
			Pattern:       pattern,
			AliasPatterns: []string{},
			Description:   "Extension " + path,
			Category:      ExtensionsCategory,
			Handler: func(context *multiplexer.Context) {
				// TODO: execute extension
			},
		})
	}
	return nil
}
