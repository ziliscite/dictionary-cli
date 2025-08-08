package engine

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/ziliscite/dictionary-cli/internal/domain"
)

type Engine struct {
	state AppState

	inputModel   *SearchModel
	loadingModel *LoadingModel

	list         list.Model
	searchResult []domain.Information
	detail       *domain.Information
}

func NewEngine() *Engine {
	inputModel := NewInputModel()
	loadingModel := NewLoadingModel()

	return &Engine{
		state:        StateSearch,
		inputModel:   inputModel,
		loadingModel: loadingModel,
	}
}

func (m *Engine) Loading() {
	m.state = StateLoading
}

func (m *Engine) Search() {
	m.state = StateSearch
}

func (m *Engine) DictionaryList() {
	m.state = StateDictionaryList
}

func (m *Engine) Detail() {
	m.state = StateDetail
}
