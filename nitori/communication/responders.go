package communication

import (
	"encoding/json"
	"errors"
	SuperVisor "git.randomchars.net/RandomChars/FreeNitori/nitori/state/supervisor"
	"syscall"
	"time"
)

func (ipc *IPC) RequestData(args []string, reply *string) error {
	for OngoingCommunication {
		time.Sleep(time.Millisecond)
	}
	OngoingCommunication = true
	defer func() {
		OngoingCommunication = false
	}()
	switch args[0] {
	case "ChatBackend":
		_ = SuperVisor.ChatBackendProcess.Signal(syscall.SIGUSR1)
		go func() {
			RequestInstructionChannel <- args[1]
		}()
		*reply = <-RequestDataChannel
	}
	return nil
}

func (ipc *IPC) RequestGuild(args []string, reply *GuildInfo) error {
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
	_ = SuperVisor.ChatBackendProcess.Signal(syscall.SIGUSR1)
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
