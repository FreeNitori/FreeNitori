package ipc

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"net/rpc"
)

var err error

// InitializeIPC dials the RPC server for any process other than the supervisor.
func InitializeIPC() error {
	if vars.ProcessType == 0 {
		return errors.New("initializing RPC client from supervisor")
	} else {
		return ipcDialClient()
	}
}

// ipcDialClient dials the RPC server.
func ipcDialClient() error {
	vars.RPCConnection, err = rpc.DialHTTP("unix", config.Config.System.Socket)
	return err
}
