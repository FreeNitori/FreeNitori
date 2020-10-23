package communication

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"git.randomchars.net/RandomChars/FreeNitori/proc/supervisor/state"
	"os"
	"strconv"
	"syscall"
)

var err error
var RequestDataChannel = make(chan string, 1)
var RequestInstructionChannel = make(chan string, 1)
var OngoingCommunication bool

type R bool

func (*R) Error(args []string, _ *int) error {
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case "ChatBackend":
		_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Error("ChatBackend has encountered a fatal error.")
		vars.ExitCode <- 1
	case "WebServer":
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
		log.Error("WebServer has encountered a fatal error.")
		vars.ExitCode <- 1
	}
	return nil
}

func (*R) Shutdown(args []string, _ *int) error {
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case "ChatBackend":
		_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by ChatBackend.")
		vars.ExitCode <- 0
	case "WebServer":
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by WebServer.")
		vars.ExitCode <- 0
	}
	return nil
}

func (*R) Restart(args []string, _ *int) error {
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case "ChatBackend":
		go func() {
			_, _ = state.ChatBackendProcess.Wait()
			state.ChatBackendProcess, err =
				os.StartProcess(config.Config.System.ChatBackend, []string{config.Config.System.ChatBackend, "-a", config.TokenOverride, "-c", config.NitoriConfPath}, &state.ProcessAttributes)
			if err != nil {
				log.Errorf("Failed to recreate chat backend process, %s", err)
				vars.ExitCode <- 1
			} else {
				log.Info("Chat backend has been restarted.")
			}
		}()
	case "WebServer":
		go func() {
			_, _ = state.WebServerProcess.Wait()
			state.WebServerProcess, err =
				os.StartProcess(config.Config.System.WebServer, []string{config.Config.System.WebServer, "-a", config.TokenOverride, "-c", config.NitoriConfPath}, &state.ProcessAttributes)
			if err != nil {
				log.Errorf("Failed to recreate web server process, %s", err)
				vars.ExitCode <- 1
			} else {
				log.Info("Web server has been restarted.")
			}
		}()
	}
	return nil
}

func (*R) FireReadyMessage(args []string, _ *int) error {
	if len(args) != 2 {
		return errors.New("invalid action")
	}
	vars.Initialized = true
	log.Infof("User: %s | ID: %s | Default Prefix: %s",
		args[0],
		args[1],
		config.Config.System.Prefix)
	log.Infof("FreeNitori is ready.")
	return nil
}

func (*R) DatabaseAction(args []string, reply *[]string) error {
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

func (*R) DatabaseActionHashmap(args []string, reply *[]map[string]string) error {
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