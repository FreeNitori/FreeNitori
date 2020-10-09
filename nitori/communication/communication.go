package communication

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
	SuperVisor "git.randomchars.net/RandomChars/FreeNitori/nitori/state/supervisor"
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
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case "ChatBackend":
		_ = SuperVisor.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Error("ChatBackend has encountered a fatal error.")
		state.ExitCode <- 1
	case "WebServer":
		_ = SuperVisor.ChatBackendProcess.Signal(syscall.SIGUSR2)
		log.Error("WebServer has encountered a fatal error.")
		state.ExitCode <- 1
	}
	return nil
}

func (*IPC) Shutdown(args []string, reply *int) error {
	reply = new(int)
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case "ChatBackend":
		_ = SuperVisor.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by ChatBackend.")
		state.ExitCode <- 0
	case "WebServer":
		_ = SuperVisor.ChatBackendProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by WebServer.")
		state.ExitCode <- 0
	}
	return nil
}

func (*IPC) Restart(args []string, reply *int) error {
	reply = new(int)
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case "ChatBackend":
		go func() {
			_, _ = SuperVisor.ChatBackendProcess.Wait()
			SuperVisor.ChatBackendProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-cb", "-a", ChatBackend.RawSession.Token, "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
			if err != nil {
				log.Errorf("Failed to recreate chat backend process, %s", err)
				state.ExitCode <- 1
			} else {
				log.Info("Chat backend has been restarted.")
			}
		}()
	case "WebServer":
		go func() {
			_, _ = SuperVisor.WebServerProcess.Wait()
			SuperVisor.WebServerProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-ws", "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
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

func (*IPC) FireReadyMessage(args []string, reply *int) error {
	reply = new(int)
	if len(args) != 2 {
		return errors.New("invalid action")
	}
	state.Initialized = true
	log.Infof("User: %s | ID: %s | Default Prefix: %s",
		args[0],
		args[1],
		config.Config.System.Prefix)
	log.Infof("FreeNitori is ready. Press Control-C to terminate.")
	return nil
}
