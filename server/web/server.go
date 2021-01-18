// Web services.
package web

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"net/http"
	"strconv"
)

var err error

func Serve() {
	<-state.DiscordReady
	log.Infof("Web server listening on %s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port)), rateLimiter.Handler(m))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to start web server, %s", err))
		state.ExitCode <- 1
	}
}
