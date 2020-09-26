package communication

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"os"
	"syscall"
)

var err error
var RequestDataChannel = make(chan string, 1)
var RequestInstructionChannel = make(chan string, 1)
var OngoingCommunication bool

type IPC bool

type GuildInfo struct {
	Name    string
	ID      string
	IconURL string
	Members []*UserInfo
}

type UserInfo struct {
	Name          string
	ID            string
	AvatarURL     string
	Discriminator string
	Bot           bool
}

func (*IPC) Error(args []string, reply *int) error {
	reply = new(int)
	switch args[0] {
	case "ChatBackend":
		_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Error("ChatBackend has encountered an error.")
		state.ExitCode <- 1
	case "WebServer":
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
		log.Error("WebServer has encountered an error.")
		state.ExitCode <- 1
	}
	return nil
}

func (*IPC) Shutdown(args []string, reply *int) error {
	reply = new(int)
	switch args[0] {
	case "ChatBackend":
		_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by ChatBackend.")
		state.ExitCode <- 0
	case "WebServer":
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by WebServer.")
		state.ExitCode <- 0
	}
	return nil
}

func (*IPC) Restart(args []string, reply *int) error {
	reply = new(int)
	switch args[0] {
	case "ChatBackend":
		go func() {
			_, _ = state.ChatBackendProcess.Wait()
			state.ChatBackendProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-cb", "-a", state.RawSession.Token, "-c", config.NitoriConfPath}, &state.ProcessAttributes)
			if err != nil {
				log.Errorf("Failed to recreate chat backend process, %s", err)
				state.ExitCode <- 1
			} else {
				log.Info("Chat backend has been restarted.")
			}
		}()
	case "WebServer":
		go func() {
			_, _ = state.WebServerProcess.Wait()
			state.WebServerProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-ws", "-c", config.NitoriConfPath}, &state.ProcessAttributes)
			if err != nil {
				log.Errorf("Failed to recreate web server process, %s", err)
				state.ExitCode <- 1
			} else {
				log.Info("Web server has been restarted.")
			}
		}()
	}
	return nil
}
