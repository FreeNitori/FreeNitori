package communication

import (
	"encoding/json"
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/ipc"
	"git.randomchars.net/RandomChars/FreeNitori/proc/supervisor/state"
	"syscall"
	"time"
)

func (*R) RequestData(args []string, reply *string) error {
	for OngoingCommunication {
		time.Sleep(time.Millisecond)
	}
	OngoingCommunication = true
	defer func() {
		OngoingCommunication = false
	}()
	switch args[0] {
	case "ChatBackend":
		_ = state.ChatBackendProcess.Signal(syscall.SIGUSR1)
		go func() {
			RequestInstructionChannel <- args[1]
		}()
		*reply = <-RequestDataChannel
	}
	return nil
}

func (*R) RequestGuild(args []string, reply *ipc.GuildInfo) error {
	for OngoingCommunication {
		time.Sleep(time.Millisecond)
	}
	OngoingCommunication = true
	defer func() {
		OngoingCommunication = false
	}()
	if len(args) == 0 {
		return errors.New("no argument was specified")
	}
	_ = state.ChatBackendProcess.Signal(syscall.SIGUSR1)
	go func() {
		RequestInstructionChannel <- "GuildInfo" + args[0]
	}()
	replyMarshalled := <-RequestDataChannel
	if replyMarshalled == "" {
		return errors.New("no guild object was returned")
	}
	err = json.Unmarshal([]byte(replyMarshalled), &reply)
	return err
}

func (*R) SignalWebServer(args []string, reply *int) error {
	args = nil
	reply = nil
	return state.WebServerProcess.Signal(syscall.SIGUSR1)
}

func (*R) ChatBackendIPCResponder(args []string, reply *string) error {
	switch args[0] {
	case "furtherInstruction":
		*reply = <-RequestInstructionChannel
	case "response":
		RequestDataChannel <- args[1]
	}
	return nil
}
