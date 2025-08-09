package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/engine"
	"os"
)

func main() {
	eng := engine.NewEngine()
	if _, err := tea.NewProgram(eng).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
