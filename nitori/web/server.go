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
		context.HTML(http.StatusOK, "web/templates/index.html", gin.H{})
	})

	// Register JSON API routes
	Engine.GET("/api", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "OK!"})
	})
	Engine.GET("/api/stats", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"total_messages":  config.GetTotalMessages(),
			"guilds_deployed": askForData("totalGuilds"),
		})
	})
	Engine.GET("/api/invite", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"invite_url": askForData("inviteURL"),
		})
	})
}
