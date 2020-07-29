package state

import (
	"github.com/bwmarrin/discordgo"
	"github.com/op/go-logging"
	"log"
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
const KappaColor = 0x3492c4

// State variables
var IPCConnection *rpc.Client
var Initialized = false
var RawSession, _ = discordgo.New()
var ShardSessions []*discordgo.Session
var Application *discordgo.Application
var Logger = logging.MustGetLogger("FreeNitori")
var format = logging.MustStringFormatter(
	logging.ColorSeq(logging.ColorGreen) + "[%{time:15:04:05.000}] %{color:reset}%{color:bold}%{level:.4s} %{color:reset}%{message}",
)
var logInfo = logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))
var logError = logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))
var ExitCode = make(chan int)
var ExecPath string
var WebServerProcess *os.Process
var ChatBackendProcess *os.Process
var ProcessAttributes = os.ProcAttr{
	Dir: ".",
	Env: os.Environ(),
	Files: []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	},
}

var EventHandlers []interface{}

func init() {
	ExecPath, err = os.Executable()
	if err != nil {
		log.Printf("Failed to get FreeNitori's executable path, %s", err)
		os.Exit(1)
	}
	logInfo.SetLevel(logging.INFO, "")
	logError.SetLevel(logging.ERROR, "")
	logging.SetBackend(logInfo, logError)
}
