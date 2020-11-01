package state

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"github.com/dgraph-io/badger/v2"
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

// Database
var Database *badger.DB
