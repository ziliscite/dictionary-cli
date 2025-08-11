package engine

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
	"strings"
)

type ExplainerDetailModel struct {
	viewport viewport.Model

	exp *domain.Explanation
}

func NewExplainerDetailModel() *ExplainerDetailModel {
	vp := viewport.New(80, 20)
	vp.Style = view.BorderStyle.PaddingRight(2)

	return &ExplainerDetailModel{
		viewport: vp,
		exp:      nil,
	}
}

func (edm *ExplainerDetailModel) Init() tea.Cmd {
	return nil
}

func (edm *ExplainerDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlQ:
			return edm, func() tea.Msg {
				return switchToExplainer{}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return edm, tea.Quit

		default:
			edm.viewport, cmd = edm.viewport.Update(msg)
			return edm, cmd
		}
	}

	return edm, cmd
}

func (edm *ExplainerDetailModel) View() string {
	fn := view.FootNoteStyle.Padding(1, 0, 2, 4).Render(
		" ↑/k up • ↓/j down • ctrl+q: back to the explainer\n",
	)

	if edm.viewport.View() == "" {
		return "No explanation" + fn
	}

	return edm.viewport.View() + fn
}

func (edm *ExplainerDetailModel) SetItem(explanation *domain.Explanation) tea.Cmd {
	edm.exp = explanation

	if edm.exp == nil {
		edm.viewport.SetContent("")
		return nil
	}

	var b strings.Builder
	core, analysis, usage := view.RenderExplainer(edm.exp)

	b.WriteString(core.Render() + "\n")
	b.WriteString(analysis.Render() + "\n")
	b.WriteString(usage.Render() + "\n")

	edm.viewport.SetContent(b.String())
	return nil
}
