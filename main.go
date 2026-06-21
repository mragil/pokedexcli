package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokemonapi"
	"pokedexcli/internal/repl"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(input *string) error
}

type ConfigApi struct {
	Next string
	Prev string
}

var commands map[string]CliCommand

const locationAreaURL = "https://pokeapi.co/api/v2/location-area/"

func initialize(pokemonapi *pokemonapi.PokemonAPI, config *ConfigApi) {
	commands = map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "Display a help message",
			Callback:    commandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Display the names of 20 location areas in the Pokemon world.",
			Callback:    commandMap(pokemonapi, config),
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display previous page the names of 20 location areas in the Pokemon world.",
			Callback:    commandMapBack(pokemonapi, config),
		},
		"explore": {
			Name:        "explore",
			Description: "list of all the Pokémon in area",
			Callback:    commandExplore(pokemonapi),
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokemonapi := pokemonapi.NewApi()
	config := ConfigApi{}

	initialize(pokemonapi, &config)

	for {
		fmt.Print("Pokedev > ")

		if start := scanner.Scan(); !start {
			break
		}

		userInput := repl.CleanInput(scanner.Text())
		if len(userInput) == 0 {
			continue
		}

		command, ok := commands[userInput[0]]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}

		if command.Name == "explore" && len(userInput) < 2 {
			fmt.Print("You must provide area to explore\n")
			continue
		}

		if command.Name != "explore" {
			err := command.Callback(nil)
			if err != nil {
				fmt.Printf("Got error: %v\n", err)
			}
			continue
		}

		err := command.Callback(&userInput[1])
		if err != nil {
			fmt.Printf("Got error: %v\n", err)
		}

	}
}

func commandExit(_ *string) error {
	_, err := fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)

	return err
}

func commandHelp(_ *string) error {
	_, err := fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, v := range commands {
		_, err := fmt.Printf("%v: %v\n", v.Name, v.Description)
		if err != nil {
			return err
		}
	}

	return err
}

func commandMap(pokemonapi *pokemonapi.PokemonAPI, config *ConfigApi) func(_ *string) error {
	return func(_ *string) error {
		url := locationAreaURL

		if config.Next != "" {
			url = config.Next
		}

		locations, next, prev, err := pokemonapi.GetLocation(url)
		if err != nil {
			return err
		}

		config.Next = next
		config.Prev = prev

		for _, location := range locations {
			fmt.Printf("%v\n", location)
		}

		return nil
	}
}

func commandMapBack(pokemonapi *pokemonapi.PokemonAPI, config *ConfigApi) func(_ *string) error {
	return func(_ *string) error {
		if config.Prev == "" {
			fmt.Println("you're on the first page")
			return nil
		}

		locations, next, prev, err := pokemonapi.GetLocation(config.Prev)
		if err != nil {
			return err
		}
		config.Next = next
		config.Prev = prev

		for _, location := range locations {
			fmt.Printf("%v\n", location)
		}

		return nil
	}
}

func commandExplore(pokemonapi *pokemonapi.PokemonAPI) func(area *string) error {
	return func(area *string) error {
		pokemons, err := pokemonapi.GetLocationAreaDetail(locationAreaURL + *area)
		if err != nil {
			return err
		}

		fmt.Printf("Exploring %v...\n", *area)
		fmt.Printf("Found Pokemon:\n")
		for _, pokemon := range pokemons {
			fmt.Printf("%v\n", pokemon)
		}

		return nil

	}
}
