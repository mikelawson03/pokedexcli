package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mikelawson03/pokedexcli/internal/api"
	"github.com/mikelawson03/pokedexcli/internal/pokedex"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

var commands map[string]cliCommand

var client = api.NewClient(15)
var pdex = pokedex.NewPokedex()

func init() {
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
			description: "Show next map location area",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show previous map location area",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Show Pokemon in area. Syntax: 'explore <area-name>'",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon. Syntax: 'catch <pokemon-name>'",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View stats of caught Pokemon. Syntax: 'inspect <pokemon-name>'",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View list of caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		cmd, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(words[1:])
		if err != nil {
			fmt.Println(err)
		}

	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	res := strings.Split(text, " ")
	return res
}

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, v := range commands {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	return nil
}

func commandMap(args []string) error {
	locationMap, err := client.GetNextLocations()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range locationMap.Results {
		fmt.Println(v.Name)
	}
	return nil
}

func commandMapb(args []string) error {
	locationMap, err := client.GetPreviousLocations()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range locationMap.Results {
		fmt.Println(v.Name)
	}
	return nil
}

func commandExplore(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("error: 'explore' requires location name")
	}

	location := args[0]
	encounters, err := client.GetEncounters(location)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range encounters.PokemonEncounters {
		fmt.Println(v.Pokemon.Name)
	}
	return nil
}

func commandCatch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("error: 'catch' requires name of Pokemon")
	}

	pokeName := args[0]
	fmt.Printf("Throwing a Pokeball at %s", pokeName)
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Print(".\n")

	pokemon, err := client.GetPokemon(pokeName)
	if err != nil {
		return err
	}
	exp := pokemon.BaseExperience

	if rand.Intn(650) >= exp {
		fmt.Println(pokeName, "was caught!")
		count, isNew := pdex.Add(pokemon)
		if isNew {
			fmt.Println(pokeName, "has been added to your Pokedex!")
		} else {
			fmt.Printf("You caught another %s. You now have %d!\n", pokeName, count)
		}
		return nil
	}
	fmt.Println(pokeName, "escaped!")
	return nil
}

func commandInspect(args []string) error {
	resp, ok := pdex.Entry[args[0]]
	if !ok {
		return fmt.Errorf("Pokemon %v not found", args[0])
	}

	pokemon := resp.Pokemon
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for k, v := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", k, v)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  -%s\n", t)
	}
	return nil
}

func commandPokedex(args []string) error {
	if len(pdex.Entry) == 0 {
		return fmt.Errorf("Pokedex is empty. Go catch a Pokemon!")
	}

	fmt.Println("Pokemon (Qty)")
	fmt.Println("-------------")
	for k, v := range pdex.Entry {
		fmt.Printf("%s (%d)\n", k, v.Count)
	}
	return nil
}
