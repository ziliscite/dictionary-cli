package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/engine"
	"net/http"
	"os"
	"time"
)

func main() {
	htc := &http.Client{
		Timeout: 5 * time.Second,
	}

	loadingModel := engine.NewLoadingModel()

	dictionaryModel := engine.NewDictionaryModel()
	detailModel := engine.NewDictionaryDetailModel()

	searchModel := engine.NewSearchModel(htc)

	deepLKey := os.Getenv("DEEPL_KEY")
	if deepLKey == "" {
		fmt.Println("DEEPL_KEY is not set")
		os.Exit(1)
	}

	translatorModel := engine.NewTranslatorModel(htc, deepLKey)

	eng := engine.NewEngine(
		searchModel,
		loadingModel,
		dictionaryModel,
		detailModel,
		translatorModel,
	)

	if _, err := tea.NewProgram(eng).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
