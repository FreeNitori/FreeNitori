package communication

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/proc/supervisor/state"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func Server() error {
	var RPCFunctions = new(R)
	_ = rpc.Register(RPCFunctions)
	rpc.HandleHTTP()
	state.SocketListener, err = net.Listen("unix", config.Config.System.Socket)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to listen on the socket, %s", err))
		os.Exit(1)
	}
	return http.Serve(state.SocketListener, nil)
}
