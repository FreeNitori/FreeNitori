package state

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"net"
	"os"
)

// RPC Server
var SocketListener net.Listener

// Service processes
var WebServerProcess *os.Process
var ChatBackendProcess *os.Process

// Process attribute
var ProcessAttributes = os.ProcAttr{
	Dir: ".",
	Env: os.Environ(),
	Files: []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	},
}

// Server arguments
var ServerArgs = []string{"-a", config.TokenOverride, "-c", config.NitoriConfPath}

// DB related
type DatabaseBackend interface {
	DBType() string
	Open(path string) error
	Close() error
	Size() int64
	GC() error
	Set(key, value string) error
	Get(key string) (string, error)
	Del(keys []string) error
	HSet(hashmap, key, value string) error
	HGet(hashmap, key string) (string, error)
	HDel(hashmap string, keys []string) error
	HGetAll(hashmap string) (map[string]string, error)
	HKeys(hashmap string) ([]string, error)
	HLen(hashmap string) (int, error)
	Iter(prefetch, includeOffset bool, offset, prefix string, handler func(key, value string) bool) error
}

var Database DatabaseBackend
