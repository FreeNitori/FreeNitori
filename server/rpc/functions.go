package rpc

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	"strconv"
)

type N bool

func (N) Version(_ []int, reply *[]string) error {
	*reply = append(*reply, state.Version)
	*reply = append(*reply, state.Revision)
	return nil
}

func (N) Shutdown(_ []int, _ *int) error {
	state.ExitCode <- 0
	return nil
}

func (N) Restart(_ []int, _ *int) error {
	state.ExitCode <- -1
	return nil
}

func (N) DatabaseAction(args []string, reply *[]string) error {
	if len(args) < 2 {
		return errors.New("invalid action")
	}
	var response = []string{""}
	switch args[0] {
	case "size":
		response[0] = strconv.Itoa(int(vars.Database.Size()))
	case "gc":
		err = vars.Database.GC()
	case "set":
		err = vars.Database.Set(args[1], args[2])
	case "get":
		response[0], err = vars.Database.Get(args[1])
	case "del":
		err = vars.Database.Del(args[1:])
	case "hset":
		err = vars.Database.HSet(args[1], args[2], args[3])
	case "hget":
		response[0], err = vars.Database.HGet(args[1], args[2])
	case "hdel":
		err = vars.Database.HDel(args[1], args[2:])
	case "hkeys":
		response, err = vars.Database.HKeys(args[1])
	case "hlen":
		var result int
		result, err = vars.Database.HLen(args[1])
		response[0] = strconv.Itoa(result)
	default:
		return errors.New("invalid operation")
	}
	*reply = response
	return err
}

func (N) DatabaseActionHashmap(args []string, reply *[]map[string]string) error {
	if len(args) < 2 {
		return errors.New("invalid action")
	}
	var response = []map[string]string{make(map[string]string)}
	switch args[0] {
	case "hgetall":
		response[0], err = vars.Database.HGetAll(args[1])
	default:
		return errors.New("invalid operation")
	}
	*reply = response
	return err
}
