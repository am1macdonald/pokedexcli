package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/am1macdonald/pokedexcli/internal/apiLink"
	"github.com/am1macdonald/pokedexcli/internal/pokecache"
	"github.com/am1macdonald/pokedexcli/internal/types/locationArea"
	"github.com/am1macdonald/pokedexcli/internal/types/locationAreaDetails"
	"github.com/am1macdonald/pokedexcli/internal/types/pokemon"
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

var commands map[string]cliCommand

var pokedex map[string]pokemon.Pokemon

func commandHelp(c *config, _ []string) error {
	fmt.Printf("%v", `
Welcome to the Pokedex!
Usage:

`)
	for _, val := range commands {
		fmt.Printf("%v: %v\n", val.name, val.description)
	}
	fmt.Println("")
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
		err = json.Unmarshal(bytes, &la)
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
	fmt.Printf("Exploring %v...\n", params[0])
	lad := locationAreaDetails.LocationAreaDetails{}
	bytes, err := getArea(params[0])
	if err != nil {
		fmt.Println("Something went wrong")
	}
	err = json.Unmarshal(bytes, &lad)
	if err != nil {
		fmt.Println("Can't marshal data!")
	}
	fmt.Println("Found pokemon:")
	for _, val := range lad.PokemonEncounters {
		fmt.Printf("  - %v\n", val.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, params []string) error {
	if len(params) <= 0 {
		fmt.Println("Pokemon name is required")
	}
	bytes, err := apiLink.FetchPokemon(params[0])
	if err != nil {
		fmt.Println("There are no pokemon with that name around")
	}
	p := pokemon.Pokemon{}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		fmt.Println("Pokemon is broken!")
	}
	fmt.Printf("Throwing a pokeball at %v...\n", params[0])
	if !p.Catch() {
		fmt.Printf("%v escaped...\n", params[0])
	} else {
		fmt.Printf("%v was caught!\n", params[0])
		pokedex[p.Name] = p
	}
	return nil
}

func commandExit(c *config, _ []string) error {
	os.Exit(0)
	return nil
}

func init() {
	cache = *pokecache.NewCache(time.Minute * 5)
	pokedex = map[string]pokemon.Pokemon{}
	commands = map[string]cliCommand{
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
			description: "List the pokemon in an area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon!",
			callback:    commandCatch,
		},
	}
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
