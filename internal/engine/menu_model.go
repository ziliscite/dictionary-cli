package engine

import (
	"fmt"
	"github.com/ziliscite/dictionary-cli/internal/view"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Choice int

const (
	Search Choice = iota
	Translate
	Explain
)

func (c Choice) String() string {
	if c < Search || c > Explain {
		return "Invalid"
	}

	return [...]string{
		"Search",
		"Translate",
		"Explain",
	}[c]
}

type MenuModel struct {
	Choice  int
	Choices []Choice
}

func NewMenuModel() *MenuModel {
	return &MenuModel{
		Choices: []Choice{
			Search, Translate, Explain,
		},
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			switch m.Choices[m.Choice] {
			case Search:
				return m, func() tea.Msg {
					return switchToSearch{}
				}

			case Translate:
				return m, func() tea.Msg {
					return switchToTranslate{}
				}

			case Explain:
				return m, func() tea.Msg {
					return switchToExplainer{}
				}
			}

		case tea.KeyDown, tea.KeyRight:
			m.Choice++
			if m.Choice > len(m.Choices)-1 {
				m.Choice = 0
			}

		case tea.KeyUp, tea.KeyLeft:
			m.Choice--
			if m.Choice < 0 {
				m.Choice = len(m.Choices) - 1
			}

		default:
			return m, nil
		}
	}

	return m, nil
}

func (m *MenuModel) View() string {
	lines := make([]string, 0, len(m.Choices))
	for i, c := range m.Choices {
		lines = append(lines, checkbox(c.String(), m.Choice == i))
	}

	choices := strings.Join(lines, "\n")

	return view.LesterViewStyle.Render(fmt.Sprintf(
		"What do you want to do?\n\n%s\n",
		choices,
	)) + view.LesterViewNoteStyle.Render(
		"esc/ctrl+c: exit • enter: choose • up/down: select\n",
	)
}

func checkbox(label string, checked bool) string {
	if checked {
		return view.DotStyle.Render("[x] " + label)
	}

	return fmt.Sprintf("[ ] %s", label)
}
