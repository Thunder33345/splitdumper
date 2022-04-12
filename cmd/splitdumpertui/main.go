package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	"os"
)

func main() {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	m := newModel(w, h)
	program := tea.NewProgram(m)

	if _, err := program.StartReturningModel(); err != nil {
		fmt.Printf("[Splitdumper] Error displaying user interface: %v\n", err)
		os.Exit(1)
	}
}
