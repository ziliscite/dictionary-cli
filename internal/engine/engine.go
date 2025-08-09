package engine

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Engine struct {
	state                AppState
	menuModel            *MenuModel
	searchModel          *SearchModel
	loadingModel         *LoadingModel
	dictionaryModel      *DictionaryModel
	detailModel          *DictionaryDetailModel
	translatorModel      *TranslatorModel
	translateDetailModel *TranslationDetailModel
}

func NewEngine(
	menuModel *MenuModel,
	searchModel *SearchModel,
	loadingModel *LoadingModel,
	dictionaryModel *DictionaryModel,
	detailModel *DictionaryDetailModel,
	translatorModel *TranslatorModel,
	translateDetailModel *TranslationDetailModel,
) *Engine {
	return &Engine{
		state:                StateMenu,
		menuModel:            menuModel,
		searchModel:          searchModel,
		loadingModel:         loadingModel,
		dictionaryModel:      dictionaryModel,
		detailModel:          detailModel,
		translatorModel:      translatorModel,
		translateDetailModel: translateDetailModel,
	}
}

func (m *Engine) Init() tea.Cmd {
	return nil
}

func (m *Engine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case switchToMenu:
		m.Menu()
		return m, nil

	case switchToDictionaryNew:
		m.DictionaryList()
		cmds = append(cmds, m.dictionaryModel.SetItems(msg.res))
		return m, tea.Batch(cmds...)

	case switchToDictionaryOld:
		m.DictionaryList()
		return m, nil

	case switchToDetail:
		m.Detail()
		cmds = append(cmds, m.detailModel.SetItem(msg.res))
		return m, tea.Batch(cmds...)

	case switchToSearch:
		m.Search()
		cmds = append(cmds, m.searchModel.Focus())
		return m, tea.Batch(cmds...)

	case switchToTranslate:
		m.Translator()
		cmds = append(cmds, m.translatorModel.Focus())
		return m, tea.Batch(cmds...)

	case switchToTranslateDetail:
		m.TranslateDetail()
		cmds = append(cmds, m.translateDetailModel.SetItem(msg.res))
		return m, tea.Batch(cmds...)

	case switchToLoading:
		m.Loading()
		cmds = append(cmds, m.loadingModel.Tick())
		return m, tea.Batch(cmds...)

	case switchToError:
		m.Menu()
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m.State(msg)
}

func updateAndAssign[T any](m *Engine, msg tea.Msg, upd func(tea.Msg) (tea.Model, tea.Cmd), set func(T)) (tea.Model, tea.Cmd) {
	mdl, cmd := upd(msg)
	if v, ok := mdl.(T); ok {
		set(v)
	}
	return m, cmd
}

func (m *Engine) State(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateMenu:
		return updateAndAssign[*MenuModel](m, msg, m.menuModel.Update, func(mm *MenuModel) {
			m.menuModel = mm
		})
	case StateLoading:
		return updateAndAssign[*LoadingModel](m, msg, m.loadingModel.Update, func(lm *LoadingModel) {
			m.loadingModel = lm
		})
	case StateSearch:
		return updateAndAssign[*SearchModel](m, msg, m.searchModel.Update, func(sm *SearchModel) {
			m.searchModel = sm
		})
	case StateDictionaryList:
		return updateAndAssign[*DictionaryModel](m, msg, m.dictionaryModel.Update, func(dm *DictionaryModel) {
			m.dictionaryModel = dm
		})
	case StateDetail:
		return updateAndAssign[*DictionaryDetailModel](m, msg, m.detailModel.Update, func(dm *DictionaryDetailModel) {
			m.detailModel = dm
		})
	case StateTranslate:
		return updateAndAssign[*TranslatorModel](m, msg, m.translatorModel.Update, func(tm *TranslatorModel) {
			m.translatorModel = tm
		})
	case StateTranslateDetail:
		return updateAndAssign[*TranslationDetailModel](m, msg, m.translateDetailModel.Update, func(tm *TranslationDetailModel) {
			m.translateDetailModel = tm
		})
	default:
		return m, nil
	}
}

func (m *Engine) View() string {
	switch m.state {
	case StateMenu:
		return m.menuModel.View()
	case StateLoading:
		return m.loadingModel.View()
	case StateSearch:
		return m.searchModel.View()
	case StateDictionaryList:
		return m.dictionaryModel.View()
	case StateDetail:
		return m.detailModel.View()
	case StateTranslate:
		return m.translatorModel.View()
	case StateTranslateDetail:
		return m.translateDetailModel.View()
	default:
		panic("unknown state")
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

func (m *Engine) Translator() {
	m.state = StateTranslate
}

func (m *Engine) TranslateDetail() {
	m.state = StateTranslateDetail
}

func (m *Engine) Menu() {
	m.state = StateMenu
}
