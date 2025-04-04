package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal"
	"os"
)

//func main() {
//	result, err := internal.Search("rain")
//	if err != nil {
//		return
//	}
//
//	internal.DisplayJisho(result)
//}

func main() {
	mod := internal.NewModel()
	if _, err := tea.NewProgram(mod).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
