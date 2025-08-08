package view

import "github.com/charmbracelet/lipgloss"

const (
	DefaultCharLimit = 56
	DefaultTextWidth = 20
	ListWidth        = 50
	ListHeight       = 14
	PaddingLeftOne   = 2
	PaddingLeftTwo   = 4

	ColorMain      = "252"
	ColorMuted     = "240"
	ColorHighlight = "170"
	ColorSpinner   = "70"
	ColorDot       = "99"
)

var (
	highlightStyle = lipgloss.NewStyle().PaddingLeft(PaddingLeftOne).Foreground(lipgloss.Color(ColorHighlight))
	normalStyle    = lipgloss.NewStyle().PaddingLeft(PaddingLeftTwo)
	wordStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMain))
	wordStyleBold  = wordStyle.Bold(true)
	mutedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted))
	mutedStyleBold = mutedStyle.Bold(true)

	dotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorDot))
)
