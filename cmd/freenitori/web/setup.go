package web

import (
	"context"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/static"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/structs"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	// Register handlers.
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	_ "git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/handlers"
	"net"
	"os"
	"syscall"
)

func Open() error {
	// Check for an existing instance if listening on unix socket
	if config.WebServer.Unix {
		if _, err := os.Stat(config.System.Socket); err != nil && !os.IsNotExist(err) {
			if _, err = net.Dial("unix", config.WebServer.Host); err != nil {
				err = syscall.Unlink(config.WebServer.Host)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("another program is listening on %s", config.WebServer.Host)
			}
		}
	}

	// Set debug mode if debug log level
	if log.GetLevel() == logrus.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Gin engine instance with recovery middleware
	router = gin.New()
	router.ForwardedByClientIP = config.WebServer.ForwardedByClientIP
	router.Use(recovery())

	// Cookie store and setup cookie middleware
	router.Use(sessions.Sessions("nitori", cookie.NewStore([]byte(config.WebServer.Secret))))

	// Enable logger if debug level
	if log.GetLevel() == logrus.DebugLevel {
		router.Use(gin.LoggerWithWriter(logger{}))
	}

	// Register templates
	if templates, err := template.ParseFS(assets, "assets/templates/*"); err != nil {
		return err
	} else {
		router.SetHTMLTemplate(templates)
	}

	// Register static
	if stat, err := os.Stat("assets/web/public"); err == nil && stat.IsDir() {
		log.Info("Serving assets from filesystem.")
		router.Use(static.ServeRoot("/", "assets/web/public"))
	} else {
		log.Info("Serving bundled assets.")
		public, err := fs.Sub(assets, "assets/public")
		if err != nil {
			return err
		}
		router.Use(static.Serve("/", static.FileSystem(http.FS(public))))
	}

	// Register error page
	router.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "error.tmpl", structs.H{
			"Title":    structs.NoSuchFileOrDirectory,
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

	// Set oauth client ID
	oauth.Conf.ClientID = state.Application.ID

	// Start server
	go serve()

	return nil
}

func Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
