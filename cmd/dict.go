package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/engine"
	"net/http"
	"os"
	"time"
	// import joho env
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		return
	}
}

func main() {
	htc := &http.Client{
		Timeout: 5 * time.Second,
	}

	menuModel := engine.NewMenuModel()
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
	translateDetailModel := engine.NewTranslationDetailModel()

	eng := engine.NewEngine(
		menuModel,
		searchModel,
		loadingModel,
		dictionaryModel,
		detailModel,
		translatorModel,
		translateDetailModel,
	)

	if _, err := tea.NewProgram(eng).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
