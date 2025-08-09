package engine

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
)

type TranslationDetailModel struct {
	tr []domain.Translation
}

func NewTranslationDetailModel() *TranslationDetailModel {
	return &TranslationDetailModel{
		tr: make([]domain.Translation, 0),
	}
}

func (ddm *TranslationDetailModel) Init() tea.Cmd {
	return nil
}

func (ddm *TranslationDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlQ:
			return ddm, func() tea.Msg {
				return switchToTranslate{}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return ddm, tea.Quit
		}
	}

	return ddm, cmd
}

func (ddm *TranslationDetailModel) View() string {
	fnt := view.FootNoteStyle.Padding(1, 0, 2, 4).Render(
		"ctrl+q: back to translation",
	)

	if ddm.tr == nil || len(ddm.tr) == 0 {
		return "" + fnt
	}

	return view.BaseViewStyle.Render(view.RenderTranslation(ddm.tr)) + fnt
}

func (ddm *TranslationDetailModel) SetItem(translations []domain.Translation) tea.Cmd {
	ddm.tr = translations
	return nil
}
