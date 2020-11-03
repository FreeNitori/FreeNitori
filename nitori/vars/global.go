// Variables containing important information.
package vars

import (
	"net/rpc"
)

// Version information
const Version = "v0.0.1-rewrite"

// Process types
const Other = -1
const Supervisor = 0
const ChatBackend = 1
const WebServer = 2

// State variables
var ProcessType int
var RPCConnection *rpc.Client
var Initialized = false
var InviteURL string
var ExitCode = make(chan int)
