package engine

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
)

type SearchModel struct {
	ti textinput.Model
	sc domain.Searcher
}

func NewInputModel() *SearchModel {
	ti := textinput.New()
	ti.Placeholder = "water"
	ti.CharLimit = 60
	ti.Width = 20
	ti.Focus()

	return &SearchModel{
		ti: ti,
	}
}

func (im *SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (im *SearchModel) View() string {
	return fmt.Sprintf(
		"What do you want to know?\n\n%s\n\n%s",
		im.ti.View(),
		"(esc to quit)",
	) + "\n"
}

func (im *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return im, tea.Batch(
				func() tea.Msg {
					return switchToLoading{}
				},
				func() tea.Msg {
					results, err := im.sc.Search(im.ti.Value())
					if err != nil {
						return switchToError{err}
					}

					return switchToDictionary{
						res: results,
					}
				},
			)

		case tea.KeyCtrlC, tea.KeyEsc:
			return im, tea.Quit
		}
	}

	im.ti, cmd = im.ti.Update(msg)
	return im, cmd
}
