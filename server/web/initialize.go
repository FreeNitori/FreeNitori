package web

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/binaries/tmpl"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	_ "git.randomchars.net/RandomChars/FreeNitori/server/web/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	ginStatic "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go/types"
	"html/template"
	"net/http"
	"strings"
)

var router *gin.Engine

type logger types.Nil

func (logger) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), err
}

func Initialize() error {

	// Set debug mode if debug log level and load certain middlewares
	if config.LogLevel == logrus.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()
	router.ForwardedByClientIP = config.Config.WebServer.ForwardedByClientIP
	router.Use(recovery())
	if config.LogLevel == logrus.DebugLevel {
		router.Use(gin.LoggerWithWriter(logger{}))
	}

	// Register templates
	templates := template.New("/")
	for _, path := range tmpl.AssetNames() {
		if strings.HasPrefix(path, "") {
			templateBin, _ := tmpl.Asset(path)
			templates, err = templates.New(path).Parse(string(templateBin))
			if err != nil {
				return errors.New("unable to parse templates")
			}
		}
	}
	router.SetHTMLTemplate(templates)

	// Register static
	router.Use(ginStatic.Serve("/", datatypes.Public()))

	// Register error page
	router.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "error.tmpl", datatypes.H{
			"Title":    datatypes.NoSuchFileOrDirectory,
			"Subtitle": "This route doesn't seem to exist.",
			"Message":  "I wonder how you got here...",
		})
	})

	// Register rate limiting middleware
	router.Use(rateMiddleware)

	for _, route := range routes.GetRoutes {
		router.GET(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.PostRoutes {
		router.POST(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.DeleteRoutes {
		router.DELETE(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.HeadRoutes {
		router.HEAD(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.OptionsRoutes {
		router.OPTIONS(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.PatchRoutes {
		router.PATCH(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.PutRoutes {
		router.PUT(route.Pattern, route.Handlers...)
	}
	for _, route := range routes.AnyRoutes {
		router.Any(route.Pattern, route.Handlers...)
	}
	return nil
}
