package engine

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
)

type DictionaryDetailModel struct {
	viewport viewport.Model

	detail *domain.Information
}

func NewDictionaryDetailModel() *DictionaryDetailModel {
	vp := viewport.New(78, 12)
	vp.Style = view.BorderStyle.PaddingRight(2)

	return &DictionaryDetailModel{
		viewport: vp,

		detail: nil,
	}
}

func (ddm *DictionaryDetailModel) Init() tea.Cmd {
	return nil
}

func (ddm *DictionaryDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			return ddm, func() tea.Msg {
				return switchToSearch{}
			}

		case tea.KeyCtrlQ:
			return ddm, func() tea.Msg {
				return switchToDictionaryOld{}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return ddm, tea.Quit

		default:
			ddm.viewport, cmd = ddm.viewport.Update(msg)
			return ddm, cmd
		}
	}

	return ddm, cmd
}

func (ddm *DictionaryDetailModel) View() string {
	return ddm.viewport.View() + view.FootNoteStyle.Padding(1, 0, 2, 4).Render(
		" ↑/k up • ↓/j down • ctrl+s: back to search • ctrl+q: back to dictionary",
	)
}

func (ddm *DictionaryDetailModel) SetItem(detail *domain.Information) tea.Cmd {
	ddm.detail = detail

	glamourRenderWidth := 78 - ddm.viewport.Style.GetHorizontalFrameSize() - 2
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
		glamour.WithColorProfile(termenv.ANSI256),
	)

	if err != nil {
		return func() tea.Msg {
			return switchToError{err}
		}
	}

	var content string
	if ddm.detail == nil {
		content = "Details not found"
	} else {
		content = view.RenderEntry(ddm.detail)
	}

	str, err := renderer.Render(content)
	if err != nil {
		return func() tea.Msg {
			return switchToError{err}
		}
	}

	ddm.viewport.SetContent(str)
	return nil
}
