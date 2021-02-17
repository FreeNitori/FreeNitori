package extension

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"github.com/robertkrimen/otto"
	"io/ioutil"
)

var jsPlugin = map[string]string{}

func registerJS(path, name string) *multiplexer.Route {
	pluginData, err := ioutil.ReadFile(path + name)
	if err != nil {
		log.Errorf("Unable to load JS extension at %s, %s", name, err)
		return nil
	}
	jsPlugin[name[:len(name)-4]] = string(pluginData)
	return &multiplexer.Route{
		Pattern:       name[:len(name)-4],
		AliasPatterns: []string{},
		Description:   "Extension loaded from " + path + name,
		Category:      ExtensionsCategory,
		Handler: func(context *multiplexer.Context) {
			err = executeJS(jsPlugin[name[:len(name)-4]], context)
			if err != nil {
				log.Errorf("Error occurred while executing extension %s, %s", name, err)
				return
			}
		},
	}
}

func executeJS(plugin string, context *multiplexer.Context) error {
	vm := otto.New()
	err = vm.Set("context", struct {
		Content     string
		SendMessage func(call otto.FunctionCall) otto.Value
	}{
		Content: context.Content,
		SendMessage: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) != 1 {
				return otto.Value{}
			}
			context.SendMessage(call.Argument(0).String())
			return otto.Value{}
		},
	})
	if err != nil {
		return err
	}
	_, err = vm.Run(plugin)
	return err
}
