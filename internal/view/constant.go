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
	HighlightStyle = lipgloss.NewStyle().PaddingLeft(PaddingLeftOne).Foreground(lipgloss.Color(ColorHighlight))
	NormalStyle    = lipgloss.NewStyle().PaddingLeft(PaddingLeftTwo)
	WordStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMain))
	WordStyleBold  = WordStyle.Bold(true)
	MutedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted))
	MutedStyleBold = MutedStyle.Bold(true)

	DotStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorDot))

	LesterViewStyle     = lipgloss.NewStyle().Padding(1, 2, 1, 0)
	LesterViewNoteStyle = MutedStyle.Padding(1, 0, 3, 0)
	BaseViewStyle       = lipgloss.NewStyle().Padding(1, 2, 1, 4)
	FootNoteStyle       = MutedStyle.Padding(1, 0, 3, 4)
)
