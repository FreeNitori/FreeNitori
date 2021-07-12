package main

import (
	"errors"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"syscall"
)

var rpcListener *net.Listener

// N represents an instance of exported stuff with all RPC receiver functions.
type N bool

// Version returns version information.
func (N) Version(_ []int, reply *[]string) error {
	*reply = append(*reply, state.Version())
	*reply = append(*reply, state.Revision())
	return nil
}

// Shutdown initiates a shutdown.
func (N) Shutdown(_ []int, _ *int) error {
	state.Exit <- 0
	return nil
}

// Restart initiates a restart.
func (N) Restart(_ []int, _ *int) error {
	state.Exit <- -1
	return nil
}

// DatabaseAction performs a database action.
func (N) DatabaseAction(args []string, reply *[]string) (err error) {
	if len(args) < 2 {
		return errors.New("invalid action")
	}
	var response = []string{""}
	switch args[0] {
	case "size":
		response[0] = strconv.Itoa(int(database.Database.Size()))
	case "gc":
		err = database.Database.GC()
	case "set":
		err = database.Database.Set(args[1], args[2])
	case "get":
		response[0], err = database.Database.Get(args[1])
	case "del":
		err = database.Database.Del(args[1:])
	case "hset":
		err = database.Database.HSet(args[1], args[2], args[3])
	case "hget":
		response[0], err = database.Database.HGet(args[1], args[2])
	case "hdel":
		err = database.Database.HDel(args[1], args[2:])
	case "hkeys":
		response, err = database.Database.HKeys(args[1])
	case "hlen":
		var result int
		result, err = database.Database.HLen(args[1])
		response[0] = strconv.Itoa(result)
	default:
		return errors.New("invalid operation")
	}
	*reply = response
	return err
}

// DatabaseActionHashmap performs a database action with hashmaps.
func (N) DatabaseActionHashmap(args []string, reply *[]map[string]string) (err error) {
	if len(args) < 2 {
		return errors.New("invalid action")
	}
	var response = []map[string]string{make(map[string]string)}
	switch args[0] {
	case "hgetall":
		response[0], err = database.Database.HGetAll(args[1])
	default:
		return errors.New("invalid operation")
	}
	*reply = response
	return err
}

func startRPC() error {
	// Check for socket in use
	if _, err := os.Stat(config.System.Socket); !os.IsNotExist(err) {
		_, err = net.Dial("unix", config.System.Socket)
		if err != nil {
			err = syscall.Unlink(config.System.Socket)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("another program is listening on %s", config.System.Socket)
		}
	}

	// Register methods
	if err := rpc.Register(new(N)); err != nil {
		return err
	}

	// Listen on and start server
	rpc.HandleHTTP()
	if listener, err := net.Listen("unix", config.System.Socket); err != nil {
		return err
	} else {
		rpcListener = &listener
		go func() { _ = http.Serve(*rpcListener, nil) }()
	}
	return nil
}
