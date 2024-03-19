package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/am1macdonald/pokedexcli/internal/apiLink"
	"github.com/am1macdonald/pokedexcli/internal/locationArea"
)

type config struct {
	next     int
	previous int
}

type cliCommand struct {
	name        string
	description string
	callback    func(config) error
}

func commandHelp(c config) error {
	fmt.Printf("%v", `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

`)
	return nil
}

func commandMap(c config) error {
	fmt.Println("Getting next map area")
	bytes, err := apiLink.FetchMap(c.next)
	if err != nil {
		fmt.Println("Something whent whrong?!")
	}
	locationArea, err := locationArea.MarshalLocationArea(bytes)
	if err != nil {
		fmt.Println("Something whent whrong?!")
	}
	fmt.Printf("%v\n", locationArea.Name)
	return nil
}

func commandMapB(c config) error {
	if c.previous == 1 {
		return errors.New("Already on first page!")
	}
	return nil
}

func commandExit(c config) error {
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
	"map": {
		name:        "map",
		description: "Display map locations",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Display previous set of map locations",
		callback:    commandMapB,
	},
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	conf := config{
		previous: 0,
		next:     1,
	}
	for {
		fmt.Print("pokedex > ")
		for sc.Scan() {
			text := sc.Text()
			if _, ok := commands[text]; ok {
				commands[text].callback(conf)
				break
			}
			fmt.Println("unknown command")
			break
		}
	}
}
