package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/engine"
	"os"
)

func main() {
	//mod := engine.NewModel()
	//if _, err := tea.NewProgram(mod).Run(); err != nil {
	//	fmt.Println("Error running program:", err)
	//	os.Exit(1)
	//}

	eng := engine.NewEngine()
	if _, err := tea.NewProgram(eng).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
