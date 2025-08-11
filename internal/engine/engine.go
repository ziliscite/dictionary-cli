package engine

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
	"reflect"
)

type Engine struct {
	state  AppState
	models map[AppState]tea.Model

	router *TransitionRouter
}

func NewEngine(
	menuModel *MenuModel,
	searchModel *SearchModel,
	loadingModel *LoadingModel,
	dictionaryModel *DictionaryModel,
	detailModel *DictionaryDetailModel,
	translatorModel *TranslatorModel,
	translateDetailModel *TranslationDetailModel,
	explainerModel *ExplainerModel,
	explainerDetailModel *ExplainerDetailModel,
) *Engine {
	models := map[AppState]tea.Model{
		StateMenu:            menuModel,
		StateLoading:         loadingModel,
		StateDictionaryList:  dictionaryModel,
		StateDetail:          detailModel,
		StateSearch:          searchModel,
		StateTranslate:       translatorModel,
		StateTranslateDetail: translateDetailModel,
		StateExplainer:       explainerModel,
		StateExplainerDetail: explainerDetailModel,
	}

	engine := &Engine{state: StateMenu, models: models, router: &TransitionRouter{
		handlers: make(map[reflect.Type]TransitionHandler),
	}}

	return engine.registerRouters()
}

func (e *Engine) registerRouters() *Engine {
	e.router.Register(switchToMenu{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		return StateMenu, nil
	})

	e.router.Register(switchToDictionaryNew{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		st := msg.(switchToDictionaryNew)
		if dm, ok := e.getModel(StateDictionaryList).(*DictionaryModel); ok {
			cmd := dm.SetItems(st.res)
			return StateDictionaryList, []tea.Cmd{cmd}
		}

		return StateDictionaryList, nil
	})

	e.router.Register(switchToDictionaryOld{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		return StateDictionaryList, nil
	})

	e.router.Register(switchToDetail{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		st := msg.(switchToDetail)
		if dm, ok := e.getModel(StateDetail).(*DictionaryDetailModel); ok {
			return StateDetail, []tea.Cmd{dm.SetItem(st.res)}
		}

		return StateDetail, nil
	})

	e.router.Register(switchToSearch{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		return StateSearch, nil
	})

	e.router.Register(switchToTranslate{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		return StateTranslate, nil
	})

	e.router.Register(switchToTranslateDetail{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		st := msg.(switchToTranslateDetail)
		if td, ok := e.getModel(StateTranslateDetail).(*TranslationDetailModel); ok {
			return StateTranslateDetail, []tea.Cmd{td.SetItem(st.res)}
		}

		return StateTranslateDetail, nil
	})

	e.router.Register(switchToLoading{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		if lm, ok := e.getModel(StateLoading).(*LoadingModel); ok {
			return StateLoading, []tea.Cmd{lm.Tick()}
		}

		return StateLoading, nil
	})

	e.router.Register(switchToError{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		err := msg.(switchToError)
		if err.err != nil {
			slog.Error(err.err.Error())
		}

		return StateMenu, nil
	})

	e.router.Register(switchToExplainer{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		return StateExplainer, nil
	})

	e.router.Register(switchToExplainerDetail{}, func(msg tea.Msg) (AppState, []tea.Cmd) {
		st := msg.(switchToExplainerDetail)
		if td, ok := e.getModel(StateExplainerDetail).(*ExplainerDetailModel); ok {
			return StateExplainerDetail, []tea.Cmd{td.SetItem(st.res)}
		}

		return StateExplainerDetail, nil
	})

	return e
}

func (e *Engine) getModel(s AppState) tea.Model {
	if m, ok := e.models[s]; ok {
		return m
	}

	return nil
}

func (e *Engine) setModel(s AppState, m tea.Model) {
	e.models[s] = m
}

func (e *Engine) Init() tea.Cmd {
	if m := e.getModel(e.state); m != nil {
		return m.Init()
	}

	return nil
}

func (e *Engine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if h, ok := e.router.Handle(msg); ok {
		nextState, cmds := h(msg)
		e.state = nextState
		if len(cmds) > 0 {
			return e, tea.Batch(cmds...)
		}

		return e, nil
	}

	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return e, tea.Quit
		}
	}

	current := e.getModel(e.state)
	if current == nil {
		return e, nil
	}

	mdl, cmd := current.Update(msg)
	e.setModel(e.state, mdl)

	return e, cmd
}

func (e *Engine) View() string {
	if m := e.getModel(e.state); m != nil {
		return m.View()
	}

	return fmt.Sprintf("Unknown state: %v", e.state)
}
