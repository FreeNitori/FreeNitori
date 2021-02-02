package rpc

import (
	"errors"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"net"
	"os"
	"syscall"
)

var err error

func Initialize() error {

	// Check for an existing instance
	if _, err := os.Stat(config.Config.System.Socket); os.IsNotExist(err) {
	} else {
		_, err := net.Dial("unix", config.Config.System.Socket)
		if err != nil {
			err = syscall.Unlink(config.Config.System.Socket)
			if err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("another program is listening on %s", config.Config.System.Socket))
		}
	}
	return nil
}
