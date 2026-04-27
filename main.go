package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	commandRegistry := newCommandRegistry()
	context := newCliContext()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		command , ok := commandRegistry[input[0]]
		if ok {
			err := command.callback(&context)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
