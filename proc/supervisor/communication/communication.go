// Inter-process communication related functions.
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

func (*R) Shutdown(args []int, _ *int) error {
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case vars.ChatBackend:
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
		_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by ChatBackend.")
		vars.ExitCode <- 0
	case vars.WebServer:
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
		_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by WebServer.")
		vars.ExitCode <- 0
	case vars.Other:
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
		_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
		log.Info("Graceful shutdown initiated by external program.")
		vars.ExitCode <- 0
	}
	return nil
}

func (*R) Restart(args []int, _ *int) error {
	if len(args) != 1 {
		return errors.New("invalid action")
	}
	switch args[0] {
	case vars.Supervisor:
		go func() {
			execPath, err := os.Executable()
			if err != nil {
				if _, err := os.Stat("bin/freenitori"); err == nil {
					execPath = "bin/freenitori"
				} else if _, err := os.Stat("build/freenitori"); err == nil {
					execPath = "build/freenitori"
				} else {
					log.Fatalf("Failed to get executable path, %s", err)
					return
				}
			}
			_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
			_, _ = state.WebServerProcess.Wait()
			_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
			_, _ = state.ChatBackendProcess.Wait()
			log.Info("Re-executing...")
			err = syscall.Exec(execPath, os.Args, os.Environ())
			if err != nil {
				log.Fatalf("Failed to re-execute, %s", err)
				vars.ExitCode <- 1
				return
			}
		}()
	case vars.ChatBackend:
		go func() {
			_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
			_, _ = state.ChatBackendProcess.Wait()
			state.ChatBackendProcess, err =
				os.StartProcess(config.Config.System.ChatBackend, append([]string{config.Config.System.ChatBackend}, state.ServerArgs...), &state.ProcessAttributes)
			if err != nil {
				log.Errorf("Failed to recreate chat backend process, %s", err)
				vars.ExitCode <- 1
			} else {
				log.Info("Chat backend has been restarted.")
			}
		}()
	case vars.WebServer:
		go func() {
			_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
			_, _ = state.WebServerProcess.Wait()
			state.WebServerProcess, err =
				os.StartProcess(config.Config.System.WebServer, append([]string{config.Config.System.WebServer}, state.ServerArgs...), &state.ProcessAttributes)
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
		response[0] = strconv.Itoa(int(state.Database.Size()))
	case "gc":
		err = state.Database.GC()
	case "set":
		err = state.Database.Set(args[1], args[2])
	case "get":
		response[0], err = state.Database.Get(args[1])
	case "del":
		err = state.Database.Del(args[1:])
	case "hset":
		err = state.Database.HSet(args[1], args[2], args[3])
	case "hget":
		response[0], err = state.Database.HGet(args[1], args[2])
	case "hdel":
		err = state.Database.HDel(args[1], args[2:])
	case "hkeys":
		response, err = state.Database.HKeys(args[1])
	case "hlen":
		var result int
		result, err = state.Database.HLen(args[1])
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
		response[0], err = state.Database.HGetAll(args[1])
	default:
		return errors.New("invalid operation")
	}
	*reply = response
	return err
}
