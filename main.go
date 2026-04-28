package main

import (
	"fmt"

	"github.com/peterh/liner"
)

func main() {
	commands := newCommandRegistry()
	context := newCliContext()

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	for {
		input, err := line.Prompt("Pokedex> ")
		if err == liner.ErrPromptAborted {
			command, _ := commands.get("exit")
			command.callback(&context)
			return
		} else if err != nil {
			// if something is wrong with our prompting mechanism we should just exit
			panic(err)
		}

		line.AppendHistory(input)
		context.input = cleanInput(input)
		if len(input) == 0 {
			fmt.Printf("You didn't type a command!")
			continue
		}

		command, ok := commands.get(context.input[0])
		if !ok {
			fmt.Printf("Invalid command %s\n", context.input[0])
			continue
		}

		err = command.callback(&context)
		if err != nil {
			fmt.Println(err)
		}
	}
}
