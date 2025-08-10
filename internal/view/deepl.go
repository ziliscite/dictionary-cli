package view

import (
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"strings"
)

func RenderTranslation(translations []domain.Translation) string {
	if len(translations) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Source: " + WordStyleBold.Underline(true).Render(translations[0].DetectedSourceLanguage) + "\n")

	for _, t := range translations {
		b.WriteString(t.Text + "\n")
	}

	return WordStyle.Render(b.String())
}
