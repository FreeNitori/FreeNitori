package communication

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"net/rpc"
)

func InitializeIPC() error {
	if state.ProcessType == 0 {
		return errors.New("initializing IPC client from supervisor")
	} else {
		return ipcDialClient()
	}
}

func ipcDialClient() error {
	state.IPCConnection, err = rpc.DialHTTP("unix", config.Config.System.Socket)
	if err != nil {
		return err
	} else {
		return nil
	}
}
