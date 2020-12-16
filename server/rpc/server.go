package rpc

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var Listener net.Listener

func Server() error {
	var Functions = new(R)
	_ = rpc.Register(Functions)
	rpc.HandleHTTP()
	Listener, err = net.Listen("unix", config.Config.System.Socket)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to listen on the socket, %s", err))
		os.Exit(1)
	}
	defer func() { _ = Listener.Close() }()
	return http.Serve(Listener, nil)
}
