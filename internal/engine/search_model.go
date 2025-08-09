package engine

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
	"net/http"
)

type SearchModel struct {
	ti textinput.Model
	sc domain.Searcher
}

func NewSearchModel(
	client *http.Client,
) *SearchModel {
	ti := textinput.New()
	ti.Placeholder = "water"
	ti.CharLimit = 60
	ti.Width = 20
	ti.Focus()

	return &SearchModel{
		ti: ti,
		sc: domain.NewSearcher(client),
	}
}

func (im *SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (im *SearchModel) View() string {
	return view.LesterViewStyle.Render(fmt.Sprintf(
		"What do you want to know?\n\n%s",
		im.ti.View(),
	)) + view.LesterViewNoteStyle.Render(
		"esc/ctrl+c: exit",
	)
}

func (im *SearchModel) searchCmd(query string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return switchToLoading{}
		},
		func() tea.Msg {
			res, err := im.sc.Search(query)
			if err != nil {
				return switchToError{err}
			}

			return switchToDictionaryNew{res: res}
		},
	)
}

func (im *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			query := im.ti.Value()
			if query == "" {
				return im, nil
			}

			im.ti.Reset()
			return im, im.searchCmd(query)

		case tea.KeyCtrlC, tea.KeyEsc:
			return im, tea.Quit

		case tea.KeyCtrlT:
			return im, func() tea.Msg {
				return switchToTranslate{}
			}
		}
	}

	im.ti, cmd = im.ti.Update(msg)
	return im, cmd
}

func (im *SearchModel) Focus() tea.Cmd {
	return im.ti.Focus()
}
