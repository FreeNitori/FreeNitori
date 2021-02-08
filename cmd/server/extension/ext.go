package extension

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"go.starlark.net/starlark"
	"io/ioutil"
	"os"
	"strings"
)

var err error

// Commands contains a map of command patterns mapped to their extension paths.
var Commands = make(map[string]string)

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
		if strings.HasSuffix(path.Name(), ".kext") {
			Commands[path.Name()[:len(path.Name())-5]] = "extensions/commands/" + path.Name()
		}
	}
	if len(Commands) > 0 {
		multiplexer.Categories = append(multiplexer.Categories, ExtensionsCategory)
	}
	return nil
}

// RegisterHandlers makes and registers event handler for each extension.
func RegisterHandlers() error {
	for pattern, path := range Commands {
		multiplexer.Router.Route(&multiplexer.Route{
			Pattern:       pattern,
			AliasPatterns: []string{},
			Description:   path,
			Category:      ExtensionsCategory,
			Handler: func(context *multiplexer.Context) {
				_, err := starlark.ExecFile(&starlark.Thread{
					Name:  pattern,
					Print: func(_ *starlark.Thread, msg string) { log.Info(msg) },
				}, path, nil, starlark.StringDict{
					"send_message": starlark.NewBuiltin("send_message", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
						if args.Len() != 1 {
							return starlark.None, errors.New("bad amount of arguments passed to send_message")
						}
						context.SendMessage(args.Index(0).String()[:len(args.Index(0).String())-1][1:])
						return starlark.None, nil
					}),
				})
				if err != nil {
					if evalErr, ok := err.(*starlark.EvalError); ok {
						log.Errorf("Extension %s encountered error while executing.", path)
						log.Error(evalErr.Backtrace())
						return
					}
					log.Errorf("Error encountered while executing extension %s, %s", path, err)
				}
			},
		})
	}
	return nil
}
