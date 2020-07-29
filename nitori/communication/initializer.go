package communication

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"net/rpc"
)

func InitializeIPC(StartChatBackend bool, StartWebServer bool) error {
	if StartChatBackend || StartWebServer {
		return ipcDialClient()
	} else {
		return errors.New("initializing IPC client from supervisor")
	}
}

func ipcDialClient() error {
	state.IPCConnection, err = rpc.DialHTTP("unix", config.SocketPath)
	if err != nil {
		return err
	} else {
		return nil
	}
}
