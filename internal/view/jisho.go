package view

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"strings"
)

func RenderEntry(entry *domain.Information) string {
	var b strings.Builder
	b.WriteString("Entry: " + wordStyleBold.Underline(true).Render(entry.Slug) + "\n")

	maxLen := max(len(entry.Japanese), len(entry.Senses))
	for i := 0; i < maxLen; i++ {
		if i < len(entry.Japanese) {
			b.WriteString("\n")
			term := entry.Japanese[i]
			if term.Word != "" {
				b.WriteString(wordStyle.Render(term.Word))
				b.WriteString(" ")
				b.WriteString(mutedStyleBold.Render("(" + term.Reading + ")"))
			} else {
				b.WriteString(wordStyle.Render(term.Reading))
			}
			b.WriteString("\n")
		}

		if i < len(entry.Senses) {
			sense := entry.Senses[i]
			b.WriteString(renderSense(sense.EnglishDefinitions, sense.PartsOfSpeech))
		}
	}

	return lipgloss.NewStyle().Render(b.String())
}

func renderSense(eng, pos []string) string {
	var b strings.Builder
	if len(pos) > 0 {
		b.WriteString(mutedStyleBold.Italic(true).Render(strings.Join(pos, ", ")))
		b.WriteString("\n")
	}

	for _, def := range eng {
		b.WriteString(dotStyle.SetString("• ").String())
		b.WriteString(wordStyle.Foreground(lipgloss.Color("252")).Render(def))
		b.WriteString("\n")
	}
	return normalStyle.MarginTop(1).Render(b.String())
}
