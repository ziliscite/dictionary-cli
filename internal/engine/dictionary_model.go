package engine

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
	"io"
	"strings"
)

type entry struct {
}

func (d entry) Height() int {
	return 1
}
func (d entry) Spacing() int {
	return 0
}
func (d entry) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d entry) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
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

type DictionaryModel struct {
	list list.Model
}

func NewDictionaryModel() *DictionaryModel {
	l := list.New(make([]list.Item, 0), entry{}, 50, 15)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = view.HighlightStyle
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select entry"),
			),
			key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "back to search"),
			),
		}
	}

	return &DictionaryModel{
		list: l,
	}
}

func (dm *DictionaryModel) Init() tea.Cmd {
	return nil
}

func (dm *DictionaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		dm.list.SetWidth(msg.Width)
		return dm, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i, ok := dm.list.SelectedItem().(item)
			switch ok {
			case true:
				return dm, func() tea.Msg {
					return switchToDetail{res: (*domain.Information)(&i)}
				}
			default:
				return dm, func() tea.Msg {
					return switchToError{err: fmt.Errorf("invalid item type")}
				}
			}

		case tea.KeyCtrlQ:
			dm.list.ResetSelected()
			return dm, func() tea.Msg {
				return switchToSearch{}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return dm, tea.Quit
		}
	}

	dm.list, cmd = dm.list.Update(msg)
	return dm, cmd
}

func (dm *DictionaryModel) SetItems(infos []domain.Information) tea.Cmd {
	items := make([]list.Item, len(infos))
	for i, info := range infos {
		items[i] = item(info)
	}

	return dm.list.SetItems(items)
}

func (dm *DictionaryModel) View() string {
	if len(dm.list.Items()) == 0 {
		return view.BaseViewStyle.Render("No items found") + view.FootNoteStyle.Render(
			"ctrl+q: back to search",
		)
	}

	return "\n" + dm.list.View()
}
