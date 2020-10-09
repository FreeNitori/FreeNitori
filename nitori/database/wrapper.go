package database

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"strconv"
)

func callDatabase(action string, data []string) (reply []string, err error) {
	body := append([]string{action}, data...)
	err = state.IPCConnection.Call("IPC.DatabaseAction", body, &reply)
	return
}

func Size() int {
	reply, _ := callDatabase("size", []string{""})
	result, _ := strconv.Atoi(reply[0])
	return result
}

func GC() error {
	_, err := callDatabase("gc", []string{""})
	return err
}

func Set(key string, value string) error {
	_, err := callDatabase("set", []string{key, value})
	return err
}

func Get(key string) (string, error) {
	reply, err := callDatabase("get", []string{key})
	return reply[0], err
}

func Del(keys []string) error {
	_, err := callDatabase("del", keys)
	return err
}
