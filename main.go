package main

import (
	"fmt"

	"github.com/yansigit/recterm/cmd"
)

func main() {
	fmt.Println("1. Record a new session")
	fmt.Println("2. Generate SVG animation")
	fmt.Print("Choose an option: ")

	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		cmd.RecordTerminal()
	case 2:
		cmd.GenerateSVG()
	default:
		fmt.Println("Invalid choice.")
	}
}
