package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	commands := newCommandRegistry()
	context := newCliContext()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		context.input = cleanInput(scanner.Text())
		command, ok := commands.get(context.input[0])
		if !ok {
			fmt.Println("invalid command")
			continue
		}

		err := command.callback(&context)
		if err != nil {
			fmt.Println(err)
		}
	}
}
