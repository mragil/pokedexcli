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
	Callback    func() error
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokemonapi := pokemonapi.NewApi()
	config := ConfigApi{}

	initialize(pokemonapi, &config)

	for {
		fmt.Print("Pokedev > ")
		if scanner.Scan() {
			userInput := repl.CleanInput(scanner.Text())
			if len(userInput) == 0 {
				continue
			}

			command, ok := commands[userInput[0]]
			if !ok {
				fmt.Print("Unknown command\n")
			} else {
				err := command.Callback()
				if err != nil {
					fmt.Printf("Got error: %v\n", err)
				}
			}
		} else {
			break
		}
	}
}

func commandExit() error {
	_, err := fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)

	return err
}

func commandHelp() error {
	_, err := fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, v := range commands {
		_, err := fmt.Printf("%v: %v\n", v.Name, v.Description)
		if err != nil {
			return err
		}
	}

	return err
}

func commandMap(pokemonapi *pokemonapi.PokemonAPI, config *ConfigApi) func() error {
	return func() error {
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

func commandMapBack(pokemonapi *pokemonapi.PokemonAPI, config *ConfigApi) func() error {
	return func() error {
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
