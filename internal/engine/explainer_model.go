package engine

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
	"log/slog"
)

type ExplainerModel struct {
	ti textinput.Model
	sc domain.Explainer
}

func NewExplainerModel(
	explainer domain.Explainer,
) *ExplainerModel {
	ti := textinput.New()
	ti.Placeholder = "私はバカな男だ"
	ti.CharLimit = 255
	ti.Width = 20
	ti.Focus()

	return &ExplainerModel{
		ti: ti,
		sc: explainer,
	}
}

func (em *ExplainerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (em *ExplainerModel) View() string {
	return view.LesterViewStyle.Render(fmt.Sprintf(
		"Insert Japanese Sentence to get the explanation: \n\n%s",
		em.ti.View(),
	)) + view.LesterViewNoteStyle.Render(
		"esc/ctrl+c: exit • ctrl+q: back to menu • enter: ask",
	)
}

func (em *ExplainerModel) askCmd(ctx context.Context, query string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return switchToLoading{}
		},
		func() tea.Msg {
			res, err := em.sc.Ask(ctx, query)
			if err != nil {
				return switchToError{err}
			}

			slog.Info("explain result", "res", res)
			return switchToExplainerDetail{res: res}
		},
	)
}

func (em *ExplainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			query := em.ti.Value()
			if query == "" {
				return em, nil
			}

			em.ti.Reset()
			return em, em.askCmd(context.Background(), query)

		case tea.KeyCtrlQ:
			em.ti.Reset()
			return em, func() tea.Msg {
				return switchToMenu{}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return em, tea.Quit
		}
	}

	em.ti, cmd = em.ti.Update(msg)
	return em, cmd
}

func (em *ExplainerModel) Focus() tea.Cmd {
	return em.ti.Focus()
}
