package engine

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
)

type DictionaryDetailModel struct {
	detail *domain.Information
}

func NewDictionaryDetailModel() *DictionaryDetailModel {
	return &DictionaryDetailModel{
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
		}
	}

	return ddm, cmd
}

func (ddm *DictionaryDetailModel) View() string {
	fnt := view.FootNoteStyle.Padding(1, 0, 2, 4).Render(
		"ctrl+s: back to search • ctrl+q: back to dictionary",
	)

	if ddm.detail == nil {
		return "" + fnt
	}

	return view.BaseViewStyle.Render(view.RenderEntry(ddm.detail)) + fnt
}

func (ddm *DictionaryDetailModel) SetItem(detail *domain.Information) tea.Cmd {
	ddm.detail = detail
	return nil
}
