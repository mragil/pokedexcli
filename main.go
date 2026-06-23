package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokemonapi"
	"pokedexcli/internal/utils"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(input *string) error
	NeedData    bool
	NeedDataMsg string
}

type ConfigApi struct {
	Next string
	Prev string
}

var commands map[string]CliCommand
var pokedex map[string]pokemonapi.Pokemon

const (
	locationAreaURL = "https://pokeapi.co/api/v2/location-area/"
	pokemonURL      = "https://pokeapi.co/api/v2/pokemon/"
)

func initialize(pokemonapi *pokemonapi.PokemonAPI, config *ConfigApi, trainer *pokemonapi.PokemonTrainer) {
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
			NeedData:    true,
			NeedDataMsg: "To use explore you must provide area name",
		},
		"catch": {
			Name:        "catch",
			Description: "Catching Pokemon and adds them to the user's Pokedex.",
			Callback:    commandCatch(pokemonapi, trainer),
			NeedData:    true,
			NeedDataMsg: "To use catch you must provide pokemon name",
		},
		"inspect": {
			Name:        "inspect",
			Description: "List all information of catched pokemon",
			Callback:    commandInspect(trainer),
			NeedData:    true,
			NeedDataMsg: "To use inspect you must provide pokemon name",
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "List all caught pokemon",
			Callback:    commandPokedex(trainer),
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokemonApi := pokemonapi.NewApi()
	config := ConfigApi{}
	trainer := pokemonapi.PokemonTrainer{
		CatchChance: 50,
		Pokedex:     make(map[string]*pokemonapi.Pokemon),
	}

	initialize(pokemonApi, &config, &trainer)

	for {
		fmt.Print("Pokedev > ")

		if start := scanner.Scan(); !start {
			break
		}

		userInput := utils.CleanInput(scanner.Text())
		if len(userInput) == 0 {
			continue
		}

		command, ok := commands[userInput[0]]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}

		if command.NeedData && len(userInput) < 2 {
			fmt.Printf("%v\n", command.NeedDataMsg)
			continue
		}

		if !command.NeedData {
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
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)

	return nil
}

func commandHelp(_ *string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, v := range commands {
		fmt.Printf("%v: %v\n", v.Name, v.Description)
	}

	return nil
}

func commandMap(pokemonApi *pokemonapi.PokemonAPI, config *ConfigApi) func(_ *string) error {
	return func(_ *string) error {
		url := locationAreaURL

		if config.Next != "" {
			url = config.Next
		}

		locations, next, prev, err := pokemonApi.GetLocation(url)
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

func commandMapBack(pokemonApi *pokemonapi.PokemonAPI, config *ConfigApi) func(_ *string) error {
	return func(_ *string) error {
		if config.Prev == "" {
			fmt.Println("you're on the first page")
			return nil
		}

		locations, next, prev, err := pokemonApi.GetLocation(config.Prev)
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

func commandExplore(pokemonApi *pokemonapi.PokemonAPI) func(area *string) error {
	return func(area *string) error {
		pokemons, err := pokemonApi.GetLocationAreaDetail(locationAreaURL + *area)
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

func commandCatch(pokemonApi *pokemonapi.PokemonAPI, trainer *pokemonapi.PokemonTrainer) func(pokemonName *string) error {
	return func(pokemonName *string) error {
		v, exist := trainer.Pokedex[*pokemonName]

		if !exist {
			pokemon, err := pokemonApi.GetPokemonDetail(pokemonURL + *pokemonName)
			if err != nil {
				return err
			}
			trainer.Pokedex[*pokemonName] = &pokemon
			v = &pokemon
		}

		if v.Chance <= 0 {
			fmt.Printf("you already caught %v!\n", v.Name)
			return nil
		}

		fmt.Printf("Throwing a Pokeball at %v...\n", v.Name)
		trainer.Pokedex[*pokemonName].Chance -= trainer.CatchChance

		if v.Chance > 0 {
			fmt.Printf("%v escaped!\n", v.Name)
			return nil
		}

		fmt.Printf("%v was caught!\n", v.Name)
		fmt.Println("You may now inspect it with the inspect command.")

		return nil
	}
}

func commandInspect(trainer *pokemonapi.PokemonTrainer) func(pokemonName *string) error {
	return func(pokemonName *string) error {
		v, exist := trainer.Pokedex[*pokemonName]

		if !exist {
			fmt.Printf("you can only inspect caught pokemon!\n")
			return nil
		}

		if v.Chance > 0 {
			fmt.Printf("you can only inspect caught pokemon!\n")
			return nil
		}

		fmt.Printf("Name: %v\n", v.Name)
		fmt.Printf("Height: %v\n", v.Height)
		fmt.Printf("Weight: %v\n", v.Weight)

		fmt.Printf("Stats:\n")
		fmt.Printf("\thp: %v\n", v.Stats.HP)
		fmt.Printf("\tattack: %v\n", v.Stats.Attack)
		fmt.Printf("\tdefense: %v\n", v.Stats.Defense)
		fmt.Printf("\tspecial-attack: %v\n", v.Stats.SpecialAttack)
		fmt.Printf("\tspecial-defense: %v\n", v.Stats.SpecialDefense)
		fmt.Printf("\tspeed: %v\n", v.Stats.Speed)

		fmt.Printf("Types:\n")

		for _, typeVal := range v.Types {
			fmt.Printf("\t- %v\n", typeVal)
		}

		return nil
	}
}

func commandPokedex(trainer *pokemonapi.PokemonTrainer) func(_ *string) error {
	return func(_ *string) error {
		if len(trainer.Pokedex) == 0 {
			fmt.Println("Your pokedex is empty. Caught some pokemon first!")
			return nil
		}

		fmt.Println("Your pokedex:")
		for _, pokemon := range trainer.Pokedex {
			fmt.Printf("\t- %v\n", pokemon.Name)
		}

		return nil
	}
}
