package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	callback    func(*config, []string) error
}

var cache pokecache.Cache

func commandHelp(c *config, _ []string) error {
	fmt.Printf("%v", `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

`)
	return nil
}

func getArea(name string) ([]byte, error) {
	bytes, ok := cache.Get(name)
	if !ok {
		b, err := apiLink.FetchMap(name)
		if err != nil {
			return nil, err
		}
		bytes = b
		cache.Add(name, b)
	}
	return bytes, nil
}

func mapThroughArea(start int) error {
	for i := start; i < (start + 20); i++ {
		bytes, err := getArea(strconv.Itoa(i))
		if err != nil {
			return err
		}
		la := locationArea.LocationArea{}
		err = la.Marshal(bytes)
		if err != nil {
			return err
		}
		fmt.Printf("%d : %v\n", i, la.Name)
	}
	return nil
}

func commandMap(c *config, _ []string) error {
	if mapThroughArea(c.next) != nil {
		fmt.Println("Something went wrong")
	}
	c.previous = c.next
	c.next += 20
	return nil
}

func commandMapB(c *config, _ []string) error {
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

func commandExplore(c *config, params []string) error {
	if len(params) <= 0 {
		fmt.Println("Area name is required!")
	}
	lad := locationArea.LocationAreaDetails{}
	bytes, err := getArea(params[0])
	if err != nil {
		fmt.Println("Something went wrong")
	}
	lad.Marshal(bytes)
	for _, val := range lad.PokemonEncounters {
		fmt.Println(val.Pokemon.Name)
	}
	return nil
}

func commandExit(c *config, _ []string) error {
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
	"explore": {
		name:        "explore",
		description: "list the pokemon in an area",
		callback:    commandExplore,
	},
}

func init() {
	cache = *pokecache.NewCache(time.Minute * 5)
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
			text := strings.Fields(sc.Text())
			if _, ok := commands[text[0]]; ok {
				commands[text[0]].callback(&conf, text[1:])
				break
			}
			fmt.Println("unknown command")
			break
		}
	}
}
