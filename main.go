package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/am1macdonald/pokedexcli/internal/apiLink"
	"github.com/am1macdonald/pokedexcli/internal/locationArea"
	"github.com/am1macdonald/pokedexcli/internal/pokecache"
)

type config struct {
	next     int
	previous int
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var cache pokecache.Cache

func commandHelp(c *config) error {
	fmt.Printf("%v", `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

`)
	return nil
}

func mapThroughArea(start int) error {
	for i := start; i < (start + 20); i++ {
		var bytes []byte
		bytes, ok := cache.Get(strconv.Itoa(i))
		if !ok {
			b, err := apiLink.FetchMap(i)
			if err != nil {
				return err
			}
			bytes = b
			cache.Add(strconv.Itoa(i), b)
		}
		locationArea, err := locationArea.MarshalLocationArea(bytes)
		if err != nil {
			return err
		}
		fmt.Printf("%d : %v\n", i, locationArea.Name)
	}
	return nil
}

func commandMap(c *config) error {
	if mapThroughArea(c.next) != nil {
		fmt.Println("Something went wrong")
	}
	c.previous = c.next
	c.next += 20
	return nil
}

func commandMapB(c *config) error {
	if c.previous < 1 {
		c.previous = 1
	}
	if mapThroughArea(c.previous) != nil {
		fmt.Println("Something went wrong")
	}
	c.next = c.previous
	c.previous -= 20
	return nil
}

func commandExit(c *config) error {
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

func init() {
	cache = *pokecache.NewCache(time.Second * 30)
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
				commands[text].callback(&conf)
				break
			}
			fmt.Println("unknown command")
			break
		}
	}
}
