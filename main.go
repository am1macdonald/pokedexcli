package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Printf("%v", `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

`)
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

var commands = map[string]cliCommand{
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		for sc.Scan() {
			text := sc.Text()
			if _, ok := commands[text]; ok {
				commands[text].callback()
				break
			}
			fmt.Println("unknown command")
			break
		}
	}
}
