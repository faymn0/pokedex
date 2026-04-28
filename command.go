package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	pokeapi "github.com/faymn0/pokedex/internal/pokeapi"
)

type cliContext struct {
	client   pokeapi.Client
	pokedex  map[string]pokeapi.PokemonResponse
	input    []string
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

type commandRegistry struct {
	data  map[string]cliCommand
	order []string
}

func newCliContext() cliContext {
	return cliContext{
		client:  pokeapi.NewClient(),
		pokedex: make(map[string]pokeapi.PokemonResponse),
	}
}

func newCommandRegistry() commandRegistry {
	commands := commandRegistry{
		data: make(map[string]cliCommand),
	}

	commands.register(cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(_ *cliContext) error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:")
			fmt.Println("")
			// in order to cyclically reference the commands object, we must make it first
			for _, key := range commands.order {
				command := commands.data[key]
				fmt.Printf("%s\t\t%s\n", command.name, command.description)
			}
			return nil
		},
	})

	commands.register(cliCommand{
		name:        "map",
		description: "Display the names of the next 20 location areas",
		callback: func(context *cliContext) error {
			return printLocationAreas(context, context.location.next)
		},
	})

	commands.register(cliCommand{
		name:        "bmap",
		description: "Display the names of the previous 20 location areas",
		callback: func(context *cliContext) error {
			return printLocationAreas(context, context.location.previous)
		},
	})

	// this is the first command that takes input arguments
	// i think ideally a robust cli would validate flags/parameters ahead of time, not inside of the command logic
	commands.register(cliCommand{
		name:        "explore",
		description: "List the pokemon in a region",
		callback: func(context *cliContext) error {
			if len(context.input) < 2 {
				return fmt.Errorf("expected a region to explore")
			}

			name := context.input[1]
			fmt.Printf("Exploring %s...\n", name)

			area, err := context.client.GetLocationArea(name)
			if err != nil {
				return err
			}

			for _, encounter := range area.PokemonEncounters {
				fmt.Printf("- %s\n", encounter.Pokemon.Name)
			}

			return nil
		},
	})

	commands.register(cliCommand{
		name:        "catch",
		description: "Catch a Pokemon!",
		callback: func(context *cliContext) error {
			name := context.input[1]
			pokemon, err := context.client.GetPokemon(name)
			if err != nil {
				// not totally accurate, we should see if it was an error because of a status code
				fmt.Println("That pokemon does not exist!")
				return nil
			}

			fmt.Printf("Throwing a Pokeball at %s...\n", name)
			if rand.Float64() < 0.5 {
				fmt.Printf("%s ran away!\n", name)
				return nil
			}

			fmt.Printf("%s was caught!\n", name)
			fmt.Println("You may now inspect it with the inspect command.")
			context.pokedex[name] = pokemon
			return nil
		},
	})

	commands.register(cliCommand{
		name:        "inspect",
		description: "Inspect a caught Pokemon",
		callback: func(context *cliContext) error {
			name := context.input[1]
			pokemon, ok := context.pokedex[name]
			if !ok {
				fmt.Printf("You haven't caught a %s yet!\n", name)
				return nil
			}

			fmt.Printf("Name: %s\n", pokemon.Name)
			fmt.Printf("Height: %d\n", pokemon.Height)
			fmt.Printf("Weight: %d\n", pokemon.Weight)
			fmt.Println("Stats:")
			for _, s := range pokemon.Stats {
				fmt.Printf("  - %s: %d\n", s.Stat.Name, s.BaseStat)
			}
			fmt.Println("Types:")
			for _, t := range pokemon.Types {
				fmt.Printf("  - %s\n", t.Type.Name)
			}

			return nil
		},
	})

	commands.register(cliCommand{
		name:        "pokedex",
		description: "View all of the Pokemon you have caught",
		callback: func(context *cliContext) error {
			fmt.Println("Your Pokedex:")
			if len(context.pokedex) == 0 {
				fmt.Println("You don't currently have any Pokemon!")
				return nil
			}
			for name := range context.pokedex {
				fmt.Printf("- %s\n", name)
			}
			return nil
		},
	})

	commands.register(cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback: func(_ *cliContext) error {
			fmt.Println("Closing the Pokedex... Goodbye!")
			os.Exit(0)
			return nil
		},
	})

	return commands
}

func (c *commandRegistry) register(command cliCommand) {
	c.data[command.name] = command
	c.order = append(c.order, command.name)
}

func (c *commandRegistry) get(name string) (cliCommand, bool) {
	command, ok := c.data[name]
	return command, ok
}

func printLocationAreas(context *cliContext, apiUrl *string) error {
	locationAreas, err := context.client.GetAllLocationAreas(apiUrl)
	if err != nil {
		return err
	}

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}
	context.location.next = locationAreas.Next
	context.location.previous = locationAreas.Previous

	return nil
}
