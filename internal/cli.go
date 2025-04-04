package internal

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

type (
	errMsg          error
	searchResultMsg []Information
	loader          bool
)

type item Information

func (i item) FilterValue() string {
	var readings string
	for _, v := range i.Japanese {
		readings += fmt.Sprintf(v.Reading, " ")
	}

	var defs string
	var poss string
	for _, v := range i.Senses {
		for _, j := range v.EnglishDefinitions {
			defs += fmt.Sprintf(j, " ")
		}
		for _, k := range v.PartsOfSpeech {
			poss += fmt.Sprintf(k, " ")
		}
	}

	return i.Slug + readings + defs + poss
}

type itemDelegate struct {
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

var itemStyle = lipgloss.NewStyle().PaddingLeft(4)

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Slug)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return itemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

type Model struct {
	err       error
	textInput textinput.Model

	searchState bool

	loading bool
	spinner spinner.Model

	list list.Model

	searchResult []Information
	pointer      int
	detail       bool
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "water"
	ti.Focus()
	ti.CharLimit = 56
	ti.Width = 20

	l := list.New(make([]list.Item, 0), itemDelegate{}, 50, 14)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "search"),
			),
		}
	}

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("70"))
	sp.Spinner = spinner.Points

	return Model{
		textInput: ti,

		searchState: true,
		loading:     false,
		spinner:     sp,

		list:    l,
		pointer: 0,
		detail:  false,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) SetItem(infos []Information) tea.Cmd {
	items := make([]list.Item, len(infos))
	for i, info := range infos {
		items[i] = item(info)
	}

	return m.list.SetItems(items)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case searchResultMsg:
		m.loading = false
		cmds = append(cmds, m.SetItem(msg))
		return m, tea.Batch(cmds...)

	case loader:
		return m, m.spinner.Tick

	case errMsg:
		m.err = msg
		m.loading = false
		return m, nil
	}

	if m.loading {
		return updateLoading(msg, m)
	}

	if m.searchState {
		return updateSearch(msg, m)
	}

	return updateDictionaryList(msg, m)
}

func updateLoading(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func updateSearch(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.searchState = false
			m.loading = true
			return m, tea.Batch(
				func() tea.Msg {
					return loader(true)
				},
				func() tea.Msg {
					results, err := SearchList(m.textInput.Value())
					if err != nil {
						return errMsg(err)
					}

					return searchResultMsg(results)
				},
			)

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func updateDictionaryList(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			m.searchState = true
			m.textInput.Focus()
			return m, nil
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return fmt.Sprintf("\n %s\n\n", m.spinner.View())
	}

	if m.searchState {
		return searchView(m)
	} else {
		return dictionaryView(m)
	}
}

func searchView(m Model) string {
	return fmt.Sprintf(
		"What do you want to know?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func dictionaryView(m Model) string {
	return "\n" + m.list.View()
}
