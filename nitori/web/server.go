package web

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var err error
var Engine *gin.Engine

func Initialize() {

	// Initialize the engine
	gin.SetMode(gin.ReleaseMode)
	Engine = gin.New()

	// Register templates
	templates := template.New("web/templates")
	for _, path := range config.AssetNames() {
		if strings.HasPrefix(path, "") {
			templateBin, _ := config.Asset(path)
			templates, err = templates.New(path).Parse(string(templateBin))
			if err != nil {
				log.Printf("Failed to parse template, %s", err)
				_ = multiplexer.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
				os.Exit(1)
			}
		}
	}
	Engine.SetHTMLTemplate(templates)

	// Register static files
	//noinspection GoUnresolvedReference
	Engine.StaticFS("/static", AssetFile())
	Engine.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "web/templates/error.html", gin.H{
			"Title":    "No such file or directory",
			"Subtitle": "This route doesn't seem to exist.",
			"Message":  "I wonder how you got here...",
		})
	})

	// Register page routes
	Engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "web/templates/index.html", nil)
	})
	Engine.GET("/guild/:gid/leaderboard", func(context *gin.Context) {
		guildInfo := fetchGuild(context.Param("gid"))
		if guildInfo == nil {
			context.HTML(http.StatusBadRequest, "web/templates/error.html", gin.H{
				"Title":    "No such file or directory",
				"Subtitle": "This guild doesn't seem to exist.",
				"Message":  "Maybe you got the wrong URL?",
			})
			return
		}
		expEnabled, err := config.ExpEnabled(guildInfo.ID)
		if err != nil {
			context.HTML(http.StatusInternalServerError, "web/templates/error.html", gin.H{
				"Title":    "Internal Server Error",
				"Subtitle": "Failed to fetch experience system enablement status.",
				"Message":  "Nitori taking a nap?",
			})
			return
		}
		if !expEnabled {
			context.HTML(http.StatusServiceUnavailable, "web/templates/error.html", gin.H{
				"Title":    "Service Unavailable",
				"Subtitle": "This feature is disabled in your guild.",
				"Message":  "Moderators don't like Nitori?",
			})
			return
		}
		context.HTML(http.StatusOK, "web/templates/leaderboard.html", gin.H{
			"GuildName": guildInfo.Name,
			"GuildIcon": guildInfo.IconURL,
		})
	})

	// Register JSON API routes
	Engine.GET("/api", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "OK!"})
	})
	Engine.GET("/api/stats", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"total_messages":  config.GetTotalMessages(),
			"guilds_deployed": fetchData("totalGuilds"),
			"nitori_version":  fetchData("version"),
		})
	})
	Engine.GET("/api/invite", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"invite_url": fetchData("inviteURL"),
		})
	})
	Engine.GET("/api/guild/:gid", func(context *gin.Context) {
		guildInfo := fetchGuild(context.Param("gid"))
		if guildInfo == nil {
			context.JSON(http.StatusNotFound, gin.H{})
			return
		}
		context.JSON(http.StatusOK, guildInfo)
	})
	Engine.GET("/api/guild/:gid/:key", func(context *gin.Context) {
		guildInfo := fetchGuild(context.Param("gid"))
		if guildInfo == nil {
			context.JSON(http.StatusNotFound, gin.H{})
			return
		}
		switch context.Param("key") {
		case "id":
			context.JSON(http.StatusOK, guildInfo.ID)
		case "name":
			context.JSON(http.StatusOK, guildInfo.Name)
		case "icon_url":
			context.JSON(http.StatusOK, guildInfo.IconURL)
		case "members":
			context.JSON(http.StatusOK, guildInfo.Members)
		default:
			context.JSON(http.StatusNotFound, gin.H{})
		}
	})
}
