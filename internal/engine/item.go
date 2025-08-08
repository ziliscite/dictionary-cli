package engine

import (
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"strings"
)

type item domain.Information

func (i item) FilterValue() string {
	var b strings.Builder

	for _, v := range i.Japanese {
		b.WriteString(v.Reading)
		b.WriteString(" ")
	}

	for _, v := range i.Senses {
		for _, j := range v.EnglishDefinitions {
			b.WriteString(j)
			b.WriteString(" ")
		}
		for _, k := range v.PartsOfSpeech {
			b.WriteString(k)
			b.WriteString(" ")
		}
	}

	return i.Slug + b.String()
}
