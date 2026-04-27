package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type commandRegistry map[string]cliCommand

func createCommandRegistry() commandRegistry {
	commands := make(map[string]cliCommand)

	commands["help"] = cliCommand {
		name: "help",
		description: "Displays a help message",
		callback: func() error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:\n")
			// in order to cyclically reference the commands object, we must make it first
			for _, command := range commands {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
			return nil
		},
	}

	commands["exit"] = cliCommand {
		name: "exit",
		description : "Exit the Pokedex",
		callback: func() error {
			fmt.Println("Closing the Pokedex... Goodbye!")
			os.Exit(0)
			return nil
		},
	}

	return commands
}


