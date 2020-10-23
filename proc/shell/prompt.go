package main

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"github.com/c-bata/go-prompt"
	"strings"
)

var rootArgs = []prompt.Suggest{
	{Text: "db", Description: "Issue database commands."},
	{Text: "exit", Description: "Exit the shell."},
}
var dbArgs = []prompt.Suggest{
	{Text: "set", Description: "Issue database command `set`"},
	{Text: "get", Description: "Issue database command `get`"},
	{Text: "del", Description: "Issue database command `del`"},
	{Text: "hset", Description: "Issue database command `hset`"},
	{Text: "hget", Description: "Issue database command `hget`"},
	{Text: "hdel", Description: "Issue database command `hdel`"},
}

func completer(document prompt.Document) []prompt.Suggest {
	fields := strings.Fields(document.Text)
	if len(fields) == 0 {
		return nil
	}
	if len(fields) == 1 && !strings.HasSuffix(document.Text, " ") {
		return prompt.FilterHasPrefix(rootArgs, document.GetWordBeforeCursor(), false)
	} else if len(fields) < 3 {
		switch fields[0] {
		case "db":
			return prompt.FilterHasPrefix(dbArgs, document.GetWordBeforeCursor(), false)
		default:
			return nil
		}
	} else {
		return nil
	}
}

func shell() {
	for {
		text := prompt.Input(fmt.Sprintf("%s[%s] > ", socketPath, vars.Version), completer,
			prompt.OptionTitle("FreeNitori Shell"))
		fields := strings.Fields(text)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "db":
			if len(fields) < 2 {
				println("Not enough arguments.")
				continue
			}
			switch fields[1] {
			case "set":
				if len(fields) != 4 {
					fmt.Println("set requires exactly 2 arguments")
					continue
				}
				err = database.Set(fields[2], fields[3])
				if err != nil {
					fmt.Printf("An error occurred while executing this command, %s", err)
					continue
				}
			case "get":
				if len(fields) != 3 {
					fmt.Println("get requires exactly 1 argument")
					continue
				}
				result, err := database.Get(fields[2])
				if err != nil {
					fmt.Printf("An error occurred while executing this command, %s", err)
					continue
				}
				println(result)
			case "del":
				if len(fields) < 3 {
					fmt.Println("del requires at least 1 argument")
					continue
				}
				err = database.Del(fields[2:])
				if err != nil {
					fmt.Printf("An error occurred while executing this command, %s", err)
					continue
				}
			case "hset":
				if len(fields) != 5 {
					fmt.Println("hset requires exactly 3 arguments")
					continue
				}
				err = database.HSet(fields[2], fields[3], fields[4])
				if err != nil {
					fmt.Printf("An error occurred while executing this command, %s", err)
					continue
				}
			case "hget":
				if len(fields) != 4 {
					fmt.Println("hget requires exactly 2 argument")
					continue
				}
				result, err := database.HGet(fields[2], fields[3])
				if err != nil {
					fmt.Printf("An error occurred while executing this command, %s", err)
					continue
				}
				println(result)
			case "hgetall":
				if len(fields) != 3 {
					fmt.Println("hget requires exactly 1 argument")
					continue
				}
				result, err := database.HGetAll(fields[2])
				if err != nil {
					fmt.Printf("An error occurred while executing this command, %s", err)
				}
				println(result)
			case "hdel":
				if len(fields) < 4 {
					fmt.Println("hdel requires at least 2 arguments")
					continue
				}
				err = database.HDel(fields[2], fields[3:]...)
				if err != nil {
					fmt.Printf("An error occurred while executing this command, %s", err)
					continue
				}
			default:
				println("Invalid argument.")
				continue
			}
		case "exit":
			vars.ExitCode <- 0
			break
		default:
			if text == "" {
				continue
			}
			println("Command not found.")
		}
	}
}
