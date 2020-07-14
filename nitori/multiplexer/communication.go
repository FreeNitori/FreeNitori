package multiplexer

import (
	"github.com/bwmarrin/discordgo"
	"github.com/op/go-logging"
	"log"
	"net/rpc"
	"os"
	"sync"
	"syscall"
)

var Logger = logging.MustGetLogger("FreeNitori")
var format = logging.MustStringFormatter(
	logging.ColorSeq(logging.ColorGreen) + "[%{time:15:04:05.000}] %{color:reset}%{color:bold}%{level:.4s} %{color:reset}%{message}",
)
var logInfo = logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))
var logError = logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))
var ExitCode = make(chan int)
var IPCConnection *rpc.Client
var RawSession, _ = discordgo.New()
var DiscordSessions []*discordgo.Session
var ExecPath string
var err error
var WebServerProcess *os.Process
var ChatBackendProcess *os.Process
var RequestDataChannel = make(chan string, 1)
var RequestInstructionChannel = make(chan string, 1)
var ProcessAttributes = os.ProcAttr{
	Dir: ".",
	Env: os.Environ(),
	Files: []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	},
}

type IPC struct {
	locker sync.RWMutex
}

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

func (*IPC) Log(args []string, reply *int) error {
	logLevel := args[0]
	logData := args[1]
	reply = new(int)
	switch logLevel {
	case "INFO":
		Logger.Info(logData)
	case "ERROR":
		Logger.Error(logData)
	case "DEBUG":
		Logger.Debug(logData)
	}
	return nil
}

func (*IPC) Error(args []string, reply *int) error {
	reply = new(int)
	switch args[0] {
	case "ChatBackend":
		_ = WebServerProcess.Signal(syscall.SIGUSR2)
		Logger.Error("ChatBackend has encountered an error.")
		ExitCode <- 1
	case "WebServer":
		_ = ChatBackendProcess.Signal(syscall.SIGUSR2)
		Logger.Error("WebServer has encountered an error.")
		ExitCode <- 1
	}
	return nil
}

func (*IPC) Shutdown(args []string, reply *int) error {
	reply = new(int)
	switch args[0] {
	case "ChatBackend":
		_ = WebServerProcess.Signal(syscall.SIGUSR2)
		Logger.Info("Graceful shutdown initiated by ChatBackend.")
		ExitCode <- 0
	case "WebServer":
		_ = ChatBackendProcess.Signal(syscall.SIGUSR2)
		Logger.Info("Graceful shutdown initiated by WebServer.")
		ExitCode <- 0
	}
	return nil
}

func (*IPC) Restart(args []string, reply *int) error {
	reply = new(int)
	switch args[0] {
	case "ChatBackend":
		go func() {
			_, _ = ChatBackendProcess.Wait()
			ChatBackendProcess, err =
				os.StartProcess(ExecPath, []string{ExecPath, "-c", "-a", RawSession.Token}, &ProcessAttributes)
			if err != nil {
				Logger.Error("Failed to recreate chat backend process, %s", err)
				ExitCode <- 1
			} else {
				Logger.Info("Chat backend has been restarted.")
			}
		}()
	case "WebServer":
		go func() {
			_, _ = WebServerProcess.Wait()
			WebServerProcess, err =
				os.StartProcess(ExecPath, []string{ExecPath, "-w"}, &ProcessAttributes)
			if err != nil {
				Logger.Error("Failed to recreate web server process, %s", err)
				ExitCode <- 1
			} else {
				Logger.Info("Web server has been restarted.")
			}
		}()
	}
	return nil
}

func (ipc *IPC) RequestData(args []string, reply *string) error {
	switch args[0] {
	case "ChatBackend":
		_ = ChatBackendProcess.Signal(syscall.SIGUSR1)
		go func() {
			RequestInstructionChannel <- args[1]
		}()
		*reply = <-RequestDataChannel
	}
	return nil
}
