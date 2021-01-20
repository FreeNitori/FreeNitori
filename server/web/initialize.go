package web

import (
	"git.randomchars.net/RandomChars/FreeNitori/binaries/static"
	"git.randomchars.net/RandomChars/FreeNitori/binaries/tmpl"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	_ "git.randomchars.net/RandomChars/FreeNitori/server/web/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	"github.com/bwmarrin/discordgo"
	"github.com/go-macaron/bindata"
	"github.com/go-macaron/oauth2"
	"github.com/go-macaron/session"
	"github.com/sirupsen/logrus"
	"go/types"
	goauth2 "golang.org/x/oauth2"
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
	macaron.Env = macaron.PROD
	m.Use(session.Sessioner())

	if config.LogLevel == logrus.DebugLevel {
		m.Use(macaron.Logger())
		m.Use(macaron.Recovery())
		macaron.Env = macaron.DEV
	} else {
		m.Use(recovery())
	}

	go func() {
		<-state.DiscordReady
		oauth2.PathCallback = "/auth/callback"
		oauth2.PathError = "/auth/error"
		oauth2.PathLogin = "/auth/login"
		oauth2.PathLogout = "/auth/logout"
		// Register Discord OAuth stuff after DiscordReady
		m.Use(oauth2.NewOAuth2Provider(&goauth2.Config{
			ClientID:     vars.Application.ID,
			ClientSecret: config.Config.Discord.ClientSecret,
			Endpoint: goauth2.Endpoint{
				AuthURL:  discordgo.EndpointOauth2 + "authorize",
				TokenURL: discordgo.EndpointOauth2 + "token",
			},
			RedirectURL: config.Config.WebServer.BaseURL + "auth/callback",
			Scopes:      []string{ScopeIdentify, ScopeGuilds},
		}))
	}()

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
