package ipc

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"net/rpc"
)

var err error

func InitializeIPC() error {
	if vars.ProcessType == 0 {
		return errors.New("initializing RPC client from supervisor")
	} else {
		return ipcDialClient()
	}
}

func ipcDialClient() error {
	vars.RPCConnection, err = rpc.DialHTTP("unix", config.Config.System.Socket)
	return err
}
