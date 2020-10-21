package state

import (
	"net/rpc"
)

// Version information
const Version = "v0.0.1-rewrite"

// Process types
const Supervisor = 0
const ChatBackend = 1
const WebServer = 2
const InteractiveConsole = 3

// State variables
var ProcessType int
var IPCConnection *rpc.Client
var Initialized = false
var InviteURL string
var ExitCode = make(chan int)
