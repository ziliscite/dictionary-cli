package engine

import (
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

type Engine struct {
	state           AppState
	searchModel     *SearchModel
	loadingModel    *LoadingModel
	dictionaryModel *DictionaryModel
	detailModel     *DictionaryDetailModel
}

func NewEngine() *Engine {
	return &Engine{
		state:           StateSearch,
		searchModel:     NewSearchModel(),
		loadingModel:    NewLoadingModel(),
		dictionaryModel: NewDictionaryModel(),
		detailModel:     NewDictionaryDetailModel(),
	}
}

func (m *Engine) Init() tea.Cmd {
	return nil
}

func (m *Engine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
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

	case switchToLoading:
		m.Loading()
		cmds = append(cmds, m.loadingModel.Tick())
		return m, tea.Batch(cmds...)

	case switchToError:
		slog.Error("something went wrong", "error", msg.err.Error())
		m.Search()
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	switch m.state {
	case StateLoading:
		mdl, cmd := m.loadingModel.Update(msg)
		if lm, ok := mdl.(*LoadingModel); ok {
			m.loadingModel = lm
		}
		return m, cmd
	case StateSearch:
		mdl, cmd := m.searchModel.Update(msg)
		if sm, ok := mdl.(*SearchModel); ok {
			m.searchModel = sm
		}
		return m, cmd
	case StateDictionaryList:
		mdl, cmd := m.dictionaryModel.Update(msg)
		if dm, ok := mdl.(*DictionaryModel); ok {
			m.dictionaryModel = dm
		}
		return m, cmd
	case StateDetail:
		mdl, cmd := m.detailModel.Update(msg)
		if dm, ok := mdl.(*DictionaryDetailModel); ok {
			m.detailModel = dm
		}
		return m, cmd
	default:
		return m, nil
	}
}

func (m *Engine) View() string {
	switch m.state {
	case StateLoading:
		return m.loadingModel.View()
	case StateSearch:
		return m.searchModel.View()
	case StateDictionaryList:
		return m.dictionaryModel.View()
	case StateDetail:
		return m.detailModel.View()
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
