package state

import (
	"net/rpc"
	"os"
)

// Version information
const Version = "v0.0.1-rewrite"

// State variables
var StartChatBackend bool
var StartWebServer bool
var IPCConnection *rpc.Client
var Initialized = false
var InviteURL string
var ExitCode = make(chan int)
var ExecPath string

func init() {
	var err error
	if ExecPath, err = os.Executable(); err != nil {
		panic(err)
	}
}
