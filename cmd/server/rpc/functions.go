package rpc

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"strconv"
)

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
	state.ExitCode <- 0
	return nil
}

// Restart initiates a restart.
func (N) Restart(_ []int, _ *int) error {
	state.ExitCode <- -1
	return nil
}

// DatabaseAction performs a database action.
func (N) DatabaseAction(args []string, reply *[]string) error {
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
func (N) DatabaseActionHashmap(args []string, reply *[]map[string]string) error {
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
