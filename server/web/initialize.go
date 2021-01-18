package web

import (
	"git.randomchars.net/RandomChars/FreeNitori/binaries/static"
	"git.randomchars.net/RandomChars/FreeNitori/binaries/tmpl"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	_ "git.randomchars.net/RandomChars/FreeNitori/server/web/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	"github.com/go-macaron/bindata"
	"github.com/sirupsen/logrus"
	"go/types"
	"gopkg.in/macaron.v1"
	"net/http"
)

var m = macaron.NewWithLogger(logger{})

type logger types.Nil

func (logger) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), err
}

func Initialize() error {
	// TODO: Forward by client IP
	m.Use(macaron.Recovery())
	if config.LogLevel == logrus.DebugLevel {
		m.Use(macaron.Logger())
	}
	m.Use(macaron.Static("static", macaron.StaticOptions{
		Prefix:      "",
		SkipLogging: true,
		IndexFile:   "index.html",
		ETag:        true,
		FileSystem: bindata.Static(bindata.Options{
			Asset:      static.Asset,
			AssetDir:   static.AssetDir,
			AssetInfo:  static.AssetInfo,
			AssetNames: static.AssetNames,
			Prefix:     ""}),
	}))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: "templates",
		TemplateFileSystem: bindata.Templates(bindata.Options{
			Asset:      tmpl.Asset,
			AssetDir:   tmpl.AssetDir,
			AssetInfo:  tmpl.AssetInfo,
			AssetNames: tmpl.AssetNames,
			Prefix:     ""}),
	}))

	m.Router.NotFound(func(context *macaron.Context) {
		context.HTML(http.StatusNotFound, "error", datatypes.H{
			"Title":    datatypes.NoSuchFileOrDirectory,
			"Subtitle": "This route doesn't seem to exist.",
			"Message":  "I wonder how you got here...",
		})
	})

	m.Router.InternalServerError(func(context *macaron.Context) {
		context.HTML(http.StatusInternalServerError, "error", datatypes.H{
			"Title":    datatypes.InternalServerError,
			"Subtitle": "Something wrong has occurred!",
			"Message":  "I wonder what happened...",
		})
	})

	for _, route := range routes.GetRoutes {
		m.Get(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.PostRoutes {
		m.Post(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.ComboRoutes {
		m.Combo(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.DeleteRoutes {
		m.Delete(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.HeadRoutes {
		m.Head(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.OptionsRoutes {
		m.Options(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.PatchRoutes {
		m.Patch(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.PutRoutes {
		m.Put(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.AnyRoutes {
		m.Any(route.Pattern, route.Handlers...)
	}
	return nil
}
