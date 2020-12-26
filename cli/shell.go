package main

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/cli/client"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/abiosoft/ishell"
)

var sh = ishell.New()

func shell() {
	sh.SetPrompt(fmt.Sprintf("FreeNitori %s(%s) > ", state.Version(), socketPath))

	sh.AddCmd(&ishell.Cmd{
		Name: "exit",
		Help: "exit the session",
		Func: func(context *ishell.Context) {
			sh.Close()
			exitCode <- 0
		},
	})

	sh.AddCmd(&ishell.Cmd{
		Name: "action",
		Help: "perform a system action",
		Func: func(context *ishell.Context) {
			choice := context.MultiChoice([]string{
				"Shutdown",
				"Restart",
			}, "Which action to perform?")
			switch choice {
			case 0:
				_ = client.Client.Call("N.Shutdown", []int{}, nil)
				context.Println("Shutdown call sent.")
				sh.Close()
				exitCode <- 0
			case 1:
				_ = client.Client.Call("N.Restart", []int{}, nil)
				context.Println("Restart call sent.")
				sh.Close()
				exitCode <- 0
			}
		},
	})

	sh.Run()
}
