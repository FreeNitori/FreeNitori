// Package web contains stuff related to web services.
package web

import (
	"embed"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/handlers"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/gin-gonic/gin"
	"go/types"
	"net"
	"net/http"
	"strconv"
	"syscall"
)

//go:embed assets
var assets embed.FS

type logger types.Nil

func (logger) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), err
}

var router *gin.Engine
var server = http.Server{}

func serve() {

	// Populate /api/info payload.
	handlers.PopulateInfoPayload()

	// Configure listener
	var listener net.Listener
	var err error
	switch config.WebServer.Unix {
	case false:
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", config.WebServer.Host, strconv.Itoa(config.WebServer.Port)))
		if err != nil {
			log.Errorf("Error listening on %s, %s", fmt.Sprintf("%s:%s", config.WebServer.Host, strconv.Itoa(config.WebServer.Port)), err)
			state.Exit <- 1
			return
		}
		log.Infof("Web server listening on %s:%s.", config.WebServer.Host, strconv.Itoa(config.WebServer.Port))
	case true:
		listener, err = net.Listen("unix", config.WebServer.Host)
		if err != nil {
			log.Errorf("Error listening on %s, %s", config.WebServer.Host, err)
			state.Exit <- 1
			return
		}

		err = syscall.Chmod(config.WebServer.Host, 0777)
		if err != nil {
			log.Errorf("Error changing permission of web server socket, %s", err)
			state.Exit <- 1
			return
		}

		log.Infof("Web server listening on unix socket %s.", config.WebServer.Host)
	}

	server.Handler = router
	err = server.Serve(listener)
	if err != nil {
		if err == http.ErrServerClosed {
			log.Info("Web server closed.")
		} else {
			log.Errorf("Web server encountered an error, %s", err)
			state.Exit <- 1
		}
	}
}
