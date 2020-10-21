package communication

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
	SuperVisor "git.randomchars.net/RandomChars/FreeNitori/nitori/state/supervisor"
	"os"
	"strconv"
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

func (*IPC) Error(args []string, _ *int) error {
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

func (*IPC) Shutdown(args []string, _ *int) error {
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

func (*IPC) Restart(args []string, _ *int) error {
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case "ChatBackend":
		go func() {
			_, _ = SuperVisor.ChatBackendProcess.Wait()
			SuperVisor.ChatBackendProcess, err =
				os.StartProcess(config.Config.System.ChatBackend, []string{config.Config.System.ChatBackend, "-a", ChatBackend.RawSession.Token, "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
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
				os.StartProcess(config.Config.System.WebServer, []string{config.Config.System.WebServer, "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
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

func (*IPC) FireReadyMessage(args []string, _ *int) error {
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

func (*IPC) DatabaseAction(args []string, reply *[]string) error {
	if len(args) < 2 {
		return errors.New("invalid action")
	}
	var response = []string{""}
	switch args[0] {
	case "size":
		response[0] = strconv.Itoa(int(size()))
	case "gc":
		err = gc()
	case "set":
		err = set(args[1], args[2])
	case "get":
		response[0], err = get(args[1])
	case "del":
		err = del(args[1:])
	case "hset":
		err = hset(args[1], args[2], args[3])
	case "hget":
		response[0], err = hget(args[1], args[2])
	case "hdel":
		err = hdel(args[1], args[2:])
	case "hkeys":
		response, err = hkeys(args[1])
	case "hlen":
		var result int
		result, err = hlen(args[1])
		response[0] = strconv.Itoa(result)
	default:
		return errors.New("invalid operation")
	}
	*reply = response
	return err
}

func (*IPC) DatabaseActionHashmap(args []string, reply *[]map[string]string) error {
	if len(args) < 2 {
		return errors.New("invalid action")
	}
	var response = []map[string]string{make(map[string]string)}
	switch args[0] {
	case "hgetall":
		response[0], err = hgetall(args[1])
	default:
		return errors.New("invalid operation")
	}
	*reply = response
	return err
}
