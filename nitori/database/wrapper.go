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

func Set(key, value string) error {
	_, err := callDatabase("set", []string{key, value})
	return err
}

func Get(key string) (string, error) {
	reply, err := callDatabase("get", []string{key})
	if len(reply) == 0 {
		return "", nil
	}
	return reply[0], err
}

func Del(keys []string) error {
	_, err := callDatabase("del", keys)
	return err
}

func HSet(hashmap, key, value string) error {
	_, err := callDatabase("hset", []string{hashmap, key, value})
	return err
}

func HGet(hashmap, key string) (string, error) {
	reply, err := callDatabase("hget", []string{hashmap, key})
	if len(reply) == 0 {
		return "", nil
	}
	return reply[0], err
}

func HDel(hashmap string, keys ...string) error {
	_, err := callDatabase("hdel", append([]string{hashmap}, keys...))
	return err
}

func HKeys(hashmap string) ([]string, error) {
	result, err := callDatabase("hkeys", []string{hashmap})
	return result, err
}

func HLen(hashmap string) (int, error) {
	result, err := callDatabase("hlen", []string{hashmap})
	var length int
	if err == nil {
		length, err = strconv.Atoi(result[0])
	}
	return length, err
}
