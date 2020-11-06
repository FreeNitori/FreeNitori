package main

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"github.com/abiosoft/ishell"
)

var shell = ishell.New()

func initShell() {
	shell.SetPrompt(fmt.Sprintf("FreeNitori %s(%s) > ", vars.Version, socketPath))

	shell.AddCmd(&ishell.Cmd{
		Name: "exit",
		Help: "exit the session",
		Func: func(context *ishell.Context) {
			shell.Close()
			vars.ExitCode <- 0
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "action",
		Help: "perform a system action",
		Func: func(context *ishell.Context) {
			choice := context.MultiChoice([]string{
				"Complete Shutdown",
				"Restart ChatBackend",
				"Restart WebServer",
			}, "Which action to perform?")
			switch choice {
			case 0:
				_ = vars.RPCConnection.Call("R.Shutdown", []int{vars.ProcessType}, nil)
				context.Println("Performing shutdown.")
				shell.Close()
				vars.ExitCode <- 0
			case 1:
				_ = vars.RPCConnection.Call("R.Restart", []int{vars.ChatBackend}, nil)
				context.Println("Restarting ChatBackend.")
			case 2:
				_ = vars.RPCConnection.Call("R.Restart", []int{vars.WebServer}, nil)
				context.Println("Restarting WebServer.")
			}
		},
	})

	shell.Run()
}
