package state

import (
	"github.com/bwmarrin/discordgo"
	"net/rpc"
	"os"
)

// Predefined error
var err error

// Version information
const Version = "v0.0.1-rewrite"

// Static messages
const InvalidArgument = "Invalid argument."
const ErrorOccurred = "Something went wrong and I am very confused! Please try again!"
const GuildOnly = "This command can only be issued from a guild."
const FeatureDisabled = "This feature is currently disabled."
const AdminOnly = "This command is only available to system administrators!"
const OperatorOnly = "This command is only available to operators!"
const PermissionDenied = "You are not allowed to issue this command!"
const MissingUser = "Specified user is not present."
const KappaColor = 0x3492c4

// State variables
var Administrator *discordgo.User
var Operator *discordgo.User
var StartChatBackend bool
var StartWebServer bool
var IPCConnection *rpc.Client
var Initialized = false
var RawSession, _ = discordgo.New()
var ShardSessions []*discordgo.Session
var Application *discordgo.Application
var InviteURL string
var ExitCode = make(chan int)
var ExecPath string
var WebServerProcess *os.Process
var ChatBackendProcess *os.Process
var EventHandlers []interface{}
var ProcessAttributes = os.ProcAttr{
	Dir: ".",
	Env: os.Environ(),
	Files: []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	},
}

func init() {
	ExecPath, err = os.Executable()
	if err != nil {
		panic(err)
	}
}
