package engine

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ziliscite/dictionary-cli/internal/view"
)

type LoadingModel struct {
	sp spinner.Model
}

func NewLoadingModel() *LoadingModel {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("70"))
	sp.Spinner = spinner.Points

	return &LoadingModel{
		sp: sp,
	}
}

func (lm *LoadingModel) Init() tea.Cmd {
	return nil
}

func (lm *LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		lm.sp, cmd = lm.sp.Update(msg)
		return lm, cmd

		//case tea.KeyMsg:
		//	switch msg.Type {
		//	case tea.KeyBackspace:
		//		return lm, func() tea.Msg {
		//			return switchToSearch{}
		//		}
		//	}
	}

	return lm, nil
}

func (lm *LoadingModel) Tick() tea.Cmd {
	return lm.sp.Tick
}

func (lm *LoadingModel) View() string {
	return view.LesterViewStyle.Render(fmt.Sprintf("Now loading %s", lm.sp.View())) + view.LesterViewNoteStyle.Render(
		"esc: exit",
	)
}
