package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal"
	"os"
)

func main() {
	mod := internal.NewModel()
	if _, err := tea.NewProgram(mod).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
