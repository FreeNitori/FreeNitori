package rpc

import (
	"net"
	"net/http"
	"net/rpc"
)

var Listener net.Listener

func Serve() {
	var Functions = new(N)
	_ = rpc.Register(Functions)
	rpc.HandleHTTP()
	_ = http.Serve(Listener, nil)
}
