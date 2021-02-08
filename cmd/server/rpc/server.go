package rpc

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"net"
	"net/http"
	"net/rpc"
)

// Listener is an instance of net.Listener used by the RPC server.
var Listener net.Listener

// Serve starts RPC server.
func Serve() {
	<-state.DiscordReady
	var Functions = new(N)
	_ = rpc.Register(Functions)
	rpc.HandleHTTP()
	Listener, err = net.Listen("unix", config.Config.System.Socket)
	if err != nil {
		log.Warnf("RPC server was unable to start, %s", err)
	}
	_ = http.Serve(Listener, nil)
}
