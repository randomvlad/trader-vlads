package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/randomvlad/trader-vlads/internal"
)

func main() {
	program := tea.NewProgram(internal.NewGame())
	_, err := program.Run()
	if err != nil {
		fmt.Printf("Program run error: %v\n", err)
		os.Exit(1)
	}
}
