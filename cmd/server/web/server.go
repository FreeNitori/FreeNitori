// Package web contains stuff related to web services.
package web

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"net"
	"net/http"
	"strconv"
	"syscall"
)

var err error

// Server contains instance of an http.Server.
var Server = http.Server{}

// Serve serves web stuff.
func Serve() {
	<-state.DiscordReady

	var listener net.Listener
	switch config.Config.WebServer.Unix {
	case false:
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port)))
		if err != nil {
			log.Errorf("Unable to listen on %s, %s", fmt.Sprintf("%s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port)), err)
			state.ExitCode <- 1
			return
		}
		log.Infof("Web server listening on %s:%s.", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port))
	case true:
		listener, err = net.Listen("unix", config.Config.WebServer.Host)
		if err != nil {
			log.Errorf("Unable to listen on %s, %s", config.Config.WebServer.Host, err)
			state.ExitCode <- 1
			return
		}

		err = syscall.Chmod(config.Config.WebServer.Host, 0777)
		if err != nil {
			log.Errorf("Unable to change permission of web server socket, %s", err)
			state.ExitCode <- 1
			return
		}

		log.Infof("Web server listening on unix socket %s.", config.Config.WebServer.Host)
	}

	Server.Handler = router
	err = Server.Serve(listener)
	if err != nil {
		if err == http.ErrServerClosed {
			log.Info("Web server closed.")
		} else {
			log.Errorf("Web server encountered an error, %s", err)
			state.ExitCode <- 1
		}
	}
}
