package engine

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
	"io"
	"net/http"
	"strings"
)

type (
	errMsg          error
	searchResultMsg []domain.Information
	loader          bool
)

type itemDelegate struct {
}

func (d itemDelegate) Height() int {
	return 1
}
func (d itemDelegate) Spacing() int {
	return 0
}
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	fn := view.NormalStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return view.HighlightStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(fmt.Sprintf("%d. %s", index+1, i.Slug)))
}

type Model struct {
	state AppState

	err       error
	textInput textinput.Model

	spinner spinner.Model

	list list.Model

	detail *domain.Information
}

func (m *Model) Loading() {
	m.state = StateLoading
}

func (m *Model) Search() {
	m.state = StateSearch
}

func (m *Model) DictionaryList() {
	m.state = StateDictionaryList
}

func (m *Model) Detail() {
	m.state = StateDetail
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
		state:     StateSearch,
		textInput: ti,
		spinner:   sp,
		list:      l,
		detail:    nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) SetItem(infos []domain.Information) tea.Cmd {
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
		m.DictionaryList()
		cmds = append(cmds, m.SetItem(msg))
		return m, tea.Batch(cmds...)

	case errMsg:
		m.err = msg
		return m, nil
	}

	switch m.state {
	case StateLoading:
		return updateLoading(msg, m)
	case StateSearch:
		return updateSearch(msg, m)
	case StateDictionaryList:
		return updateDictionaryList(msg, m)
	case StateDetail:
		return updateDetail(msg, m)
	default:
		return m, nil
	}
}

func updateLoading(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, m.spinner.Tick
}

var sc = domain.NewSearcher(http.DefaultClient)

func updateSearch(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Loading()
			return m, tea.Batch(
				func() tea.Msg {
					return loader(true)
				},
				func() tea.Msg {
					results, err := sc.Search(m.textInput.Value())
					if err != nil {
						return errMsg(err)
					}

					return searchResultMsg(results)
				},
				func() tea.Msg {
					return loader(false)
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

			m.detail = (*domain.Information)(&i)
			m.Detail()
			return m, nil

		case tea.KeyCtrlS:
			m.textInput.Focus()
			m.list.ResetSelected()
			m.Search()
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
			m.DictionaryList()
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlS:
			m.Search()
			m.textInput.Focus()
			return m, nil
		}
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case StateLoading:
		return fmt.Sprintf("\n %s\n\n", m.spinner.View())
	case StateSearch:
		return searchView(m)
	case StateDictionaryList:
		return dictionaryView(m)
	case StateDetail:
		return detailView(m)
	default:
		panic("unknown state")
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
	if len(m.list.Items()) == 0 {
		return view.BaseViewStyle.Render("No items found") + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1, 0, 3, 4).Render(
			"ctrl+s: back to search",
		)
	}

	return "\n" + m.list.View()
}

func detailView(m Model) string {
	return view.BaseViewStyle.Render(m.renderContent()) + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1, 0, 2, 4).Render(
		"ctrl+s: back to search • backspace: back to dictionary",
	)
}

func (m Model) renderContent() string {
	if m.detail == nil {
		return ""
	}

	return view.RenderEntry(m.detail)
}
