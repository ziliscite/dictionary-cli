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

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Slug)

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render("> " + strings.Join(s, " "))
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

	detailState bool
	detail      *Information
	ready       bool
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "water"
	ti.Focus()
	ti.CharLimit = 56
	ti.Width = 20

	l := list.New(make([]list.Item, 0), itemDelegate{}, 50, 14)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
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

		list:   l,
		detail: nil,
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

	if m.detailState {
		return updateDetail(msg, m)
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
		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(item)
			if !ok {
				return m, func() tea.Msg {
					return errMsg(fmt.Errorf("item is not a dictionary"))
				}
			}

			m.detail = (*Information)(&i)
			m.detailState = true
			return m, nil

		case tea.KeyCtrlS:
			m.searchState = true
			m.textInput.Focus()
			m.list.ResetSelected()
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateDetail(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:
			m.detailState = false
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlS:
			m.detailState = false
			m.searchState = true
			m.textInput.Focus()
			return m, nil
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return fmt.Sprintf("\n %s\n\n", m.spinner.View())
	}

	if m.detailState {
		return detailView(m)
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

func detailView(m Model) string {
	return lipgloss.NewStyle().Padding(1, 2, 0, 4).Render(m.renderContent()) + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1, 0, 2, 4).Render(
		"ctrl+s: back to search • backspace: back to dictionary",
	)
}

func (m Model) renderContent() string {
	if m.detail == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(renderEntry(m.detail))
	return b.String()
}

var wordStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))

func renderEntry(entry *Information) string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("36")).Underline(true).Render(entry.Slug))
	b.WriteString("\n\n")

	for _, term := range entry.Japanese {
		if term.Word != "" {
			b.WriteString(wordStyle.Render(term.Word))
			b.WriteString(" ")
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("(" + term.Reading + ")"))
		} else {
			b.WriteString(wordStyle.Render(term.Reading))
		}
		b.WriteString("\n")
	}

	for _, sense := range entry.Senses {
		b.WriteString(renderSense(sense.EnglishDefinitions, sense.PartsOfSpeech))
	}

	return lipgloss.NewStyle().Render(b.String())
}

func renderSense(eng, pos []string) string {
	var b strings.Builder
	if len(pos) > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true).Render(strings.Join(pos, ", ")))
		b.WriteString("\n")
	}

	for _, def := range eng {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("99")).SetString("•").String() + " ")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(def))
		b.WriteString("\n")
	}
	return lipgloss.NewStyle().PaddingLeft(2).MarginTop(1).Render(b.String())
}
