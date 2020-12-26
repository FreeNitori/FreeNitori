// Web services.
package web

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/gin-gonic/gin"
	"strconv"
)

var err error
var Engine *gin.Engine

func Serve() {
	<-state.DiscordReady
	log.Infof("Web server listening on %s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port))
	err = Engine.Run(fmt.Sprintf("%s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port)))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to start web server, %s", err))
		state.ExitCode <- 1
	}
}
