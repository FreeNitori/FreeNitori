package database

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"strconv"
)

// callDatabase provides a wrapper around the RPC function DatabaseAction and performs a database call.
func callDatabase(action string, data []string) (reply []string, err error) {
	body := append([]string{action}, data...)
	err = vars.RPCConnection.Call("R.DatabaseAction", body, &reply)
	return
}

// callDatabaseHashmap provides a wrapper around the RPC function DatabaseActionHashmap and performs a database call that contains hashmaps.
func callDatabaseHashmap(action string, data []string) (reply []map[string]string, err error) {
	body := append([]string{action}, data...)
	err = vars.RPCConnection.Call("R.DatabaseActionHashmap", body, &reply)
	return
}

// Size returns the size.
func Size() int {
	reply, _ := callDatabase("size", []string{""})
	result, _ := strconv.Atoi(reply[0])
	return result
}

// GC starts the garbage collection.
func GC() error {
	_, err := callDatabase("gc", []string{""})
	return err
}

// Set sets the value of a key.
func Set(key, value string) error {
	_, err := callDatabase("set", []string{key, value})
	return err
}

// Get gets the value of a key.
func Get(key string) (string, error) {
	reply, err := callDatabase("get", []string{key})
	if len(reply) == 0 {
		return "", nil
	}
	return reply[0], err
}

// Del deletes a key.
func Del(keys []string) error {
	_, err := callDatabase("del", keys)
	return err
}

// HSet sets the value of a key in a hashmap.
func HSet(hashmap, key, value string) error {
	_, err := callDatabase("hset", []string{hashmap, key, value})
	return err
}

// HGet gets the value of a key in a hashmap.
func HGet(hashmap, key string) (string, error) {
	reply, err := callDatabase("hget", []string{hashmap, key})
	if len(reply) == 0 {
		return "", nil
	}
	return reply[0], err
}

// HGetAll gets an entire hashmap.
func HGetAll(hashmap string) (map[string]string, error) {
	reply, err := callDatabaseHashmap("hgetall", []string{hashmap})
	if len(reply) == 0 {
		return nil, nil
	}
	return reply[0], err
}

// HDel deletes a hashmap or one of the keys in it.
func HDel(hashmap string, keys ...string) error {
	_, err := callDatabase("hdel", append([]string{hashmap}, keys...))
	return err
}

// HKeys gets all keys of a hashmap.
func HKeys(hashmap string) ([]string, error) {
	result, err := callDatabase("hkeys", []string{hashmap})
	return result, err
}

// HLen measures the length of a hashmap.
func HLen(hashmap string) (int, error) {
	result, err := callDatabase("hlen", []string{hashmap})
	var length int
	if err == nil {
		length, err = strconv.Atoi(result[0])
	}
	return length, err
}
