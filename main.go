package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commandRegistry := createCommandRegistry()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		command , ok := commandRegistry[input[0]]
		if ok {
			err := command.callback()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
