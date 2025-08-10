package view

import (
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"strconv"
	"strings"
)

func RenderEntry(entry *domain.Information) string {
	var b strings.Builder
	b.WriteString("Entry: **" + entry.Slug + "**")
	if len(entry.JLPT) > 0 {
		b.WriteString(" | " + strings.Join(entry.JLPT, ", "))
	}
	b.WriteString("\n\n")

	maxLen := max(len(entry.Japanese), len(entry.Senses))
	for i := 0; i < maxLen; i++ {
		if i < len(entry.Japanese) {
			b.WriteString(strconv.Itoa(i+1) + ". ")
			term := entry.Japanese[i]
			if term.Word != "" {
				b.WriteString("**" + term.Word + "** ")
				b.WriteString("_(" + term.Reading + ")_\n")
			} else {
				b.WriteString("*" + term.Reading + "*\n")
			}
			b.WriteString("\n")
		}

		if i < len(entry.Senses) {
			sense := entry.Senses[i]
			b.WriteString(renderSensePlain(sense.EnglishDefinitions, sense.PartsOfSpeech))
		}
	}

	return b.String()
}

func renderSensePlain(eng, pos []string) string {
	var b strings.Builder
	if len(pos) > 0 {
		b.WriteString("_(" + strings.Join(pos, ", ") + ")_\n")
	}

	for _, def := range eng {
		b.WriteString("- " + def + "\n")
	}
	b.WriteString("\n")
	return b.String()
}
