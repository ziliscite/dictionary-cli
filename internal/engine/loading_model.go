package engine

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	}

	return lm, lm.sp.Tick
}

func (lm *LoadingModel) View() string {
	return fmt.Sprintf("\n %s\n\n", lm.sp.View())
}
