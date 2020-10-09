package SuperVisor

import (
	"github.com/dgraph-io/badger/v2"
	"net"
	"os"
)

// IPC Server
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

// Database
var Database *badger.DB
