package engine

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
)

type ExplainerDetailModel struct {
	exp *domain.Explanation
}

func (e *ExplainerDetailModel) Init() tea.Cmd {
	//TODO implement me
	panic("implement me")
}

func (e *ExplainerDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//TODO implement me
	panic("implement me")
}

func (e *ExplainerDetailModel) View() string {
	return ""
}

func (e *ExplainerDetailModel) SetItem(explanation *domain.Explanation) tea.Cmd {
	e.exp = explanation
	return nil
}
