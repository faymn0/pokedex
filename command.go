package main

import (
	"fmt"
	"os"

	pokeapi "github.com/faymn0/pokedex/internal"
)

type cliContext struct {
	client   pokeapi.Client
	location struct {
		next     *string
		previous *string
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cliContext) error
}

type commandRegistry map[string]cliCommand

func newCliContext() cliContext {
	return cliContext{
		client: pokeapi.NewClient(),
	}
}

func newCommandRegistry() commandRegistry {
	commands := make(map[string]cliCommand)

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(_ *cliContext) error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:")
			fmt.Println("")
			// in order to cyclically reference the commands object, we must make it first
			for _, command := range commands {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
			return nil
		},
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "Display the names of the next 20 location areas",
		callback: func(context *cliContext) error {
			locationAreas, err := context.client.GetLocationArea(context.location.next)
			if err != nil {
				return err
			}

			for _, area := range locationAreas.Results {
				fmt.Println(area.Name)
			}
			return nil
		},
	}

	commands["bmap"] = cliCommand{
		name:        "bmap",
		description: "Display the names of the previous 20 location areas",
		callback: func(context *cliContext) error {
			locationAreas, err := context.client.GetLocationArea(context.location.previous)
			if err != nil {
				return err
			}
			for _, area := range locationAreas.Results {
				fmt.Println(area.Name)
			}
			return nil
		},
	}

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback: func(_ *cliContext) error {
			fmt.Println("Closing the Pokedex... Goodbye!")
			os.Exit(0)
			return nil
		},
	}

	return commands
}

