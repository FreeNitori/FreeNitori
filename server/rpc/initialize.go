package rpc

import (
	"errors"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"net/rpc/jsonrpc"
)

var err error

func Initialize() error {
	_, err = jsonrpc.Dial("unix", config.Config.System.Socket)
	if err == nil {
		return errors.New(fmt.Sprintf("another program is already listening on %s", config.Config.System.Socket))
	}
	return nil
}
