package engine

import (
	tea "github.com/charmbracelet/bubbletea"
	"reflect"
)

type TransitionHandler func(msg tea.Msg) (nextState AppState, cmds []tea.Cmd)

type TransitionRouter struct {
	handlers map[reflect.Type]TransitionHandler
}

func (r *TransitionRouter) Register(msgSample tea.Msg, h TransitionHandler) {
	r.handlers[reflect.TypeOf(msgSample)] = h
}

func (r *TransitionRouter) Handle(msg tea.Msg) (TransitionHandler, bool) {
	h, ok := r.handlers[reflect.TypeOf(msg)]
	return h, ok
}
