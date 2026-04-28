package main

import "strings"

func cleanInput(text string) []string {
	var result []string
	words := strings.Split(text, " ")
	for _, w := range words {
		f := strings.Trim(strings.ToLower(w), " ")
		if len(f) > 0 {
			result = append(result, f)
		}
	}
	return result
}
