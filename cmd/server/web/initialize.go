package web

import (
	"errors"
	"git.randomchars.net/FreeNitori/FreeNitori/binaries/public"
	"git.randomchars.net/FreeNitori/FreeNitori/binaries/tmpl"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	log "git.randomchars.net/FreeNitori/Log"

	// Register handlers.
	_ "git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/handlers"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/static"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go/types"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var router *gin.Engine

type logger types.Nil

func (logger) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), err
}

// Initialize early initializes web services.
func Initialize() error {

	router = ginSetup()

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
	var serveFileSystem static.ServeFileSystem
	serveFileSystem = static.FileSystem(public.AssetFile())
	if stat, err := os.Stat("assets/web/public"); err == nil {
		if stat.IsDir() {
			serveFileSystem = static.LocalFile("assets/web/public", false)
		}
	}
	router.Use(static.Serve("/", serveFileSystem))

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

	// Register routes
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

func ginSetup() *gin.Engine {
	// Set debug mode if debug log level and load certain middlewares
	if log.GetLevel() == logrus.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.ForwardedByClientIP = config.Config.WebServer.ForwardedByClientIP
	engine.Use(recovery())

	store := cookie.NewStore([]byte(config.Config.WebServer.Secret))
	engine.Use(sessions.Sessions("nitori", store))

	if log.GetLevel() == logrus.DebugLevel {
		engine.Use(gin.LoggerWithWriter(logger{}))
	}
	return engine
}
