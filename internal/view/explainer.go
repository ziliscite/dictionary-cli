package view

import (
	"fmt"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"strings"
)

// RenderExplainer will split the explanation into 3 parts: Core, Analysis, and Usage
// It will return 3 of them in the form of a struct, which will have their own render functions.
func RenderExplainer(entry domain.Explanation) (*Core, *Analysis, *Usage) {
	core := &Core{
		Original:            entry.Original,
		Kana:                entry.Kana,
		Romaji:              entry.Romaji,
		LiteralTranslation:  entry.LiteralTranslation,
		NaturalTranslations: entry.NaturalTranslations,
		Confidence:          entry.Confidence,
	}

	analysis := &Analysis{
		WordByWord:    entry.WordByWord,
		GrammarPoints: entry.GrammarPoints,
	}

	usage := &Usage{
		NuanceAndRegister:          entry.NuanceAndRegister,
		CommonErrors:               entry.CommonErrors,
		ParaphrasesAndAlternatives: entry.ParaphrasesAndAlternatives,
	}

	return core, analysis, usage
}

type Core struct {
	Original            string   `json:"original"`
	Kana                string   `json:"kana"`
	Romaji              string   `json:"romaji"`
	LiteralTranslation  string   `json:"literal_translation"`
	NaturalTranslations []string `json:"natural_translations"`
	Confidence          string   `json:"confidence"`
}

func (c Core) Render() string {
	var b strings.Builder

	b.WriteString(WordStyleBold.Render("Sentence: " + c.Original + "( " + c.Kana + " )"))
	b.WriteString("\n")
	b.WriteString(WordStyle.Render("Romaji: " + c.Romaji))
	b.WriteString(WordStyle.Render("Translations: "))

	b.WriteString(DotStyle.Render(">") + " Literal:\n")
	b.WriteString(WordStyle.Render(c.LiteralTranslation))
	b.WriteString(DotStyle.Render(">") + " Natural:\n")
	b.WriteString(WordStyle.Render(strings.Join(c.NaturalTranslations, "\n")))

	b.WriteString(WordStyle.Render("Confidence: " + c.Confidence))

	return b.String()
}

type Analysis struct {
	WordByWord []struct {
		Token   string `json:"token"`
		Reading string `json:"reading"`
		Pos     string `json:"pos"`
		Meaning string `json:"meaning"`
	} `json:"word_by_word"`
	GrammarPoints []struct {
		Point           string   `json:"point"`
		Explanation     string   `json:"explanation"`
		SimilarExamples []string `json:"similar_examples"`
	} `json:"grammar_points"`
}

func (a Analysis) Render() string {
	var b strings.Builder

	b.WriteString(WordStyleBold.Render("Gloss Analysis: ") + "\n")

	for _, v := range a.WordByWord {
		b.WriteString(WordStyleBold.Italic(true).Render(v.Token + " "))
		b.WriteString(MutedStyle.Render("("+v.Reading+")") + DotStyle.Render("→ "))
		b.WriteString(WordStyle.Render(v.Pos + "(" + v.Meaning + ")" + "\n"))
	}

	b.WriteString("\n")
	b.WriteString(WordStyleBold.Render("Grammar Points: ") + "\n")

	for i, p := range a.GrammarPoints {
		b.WriteString(DotStyle.Render(fmt.Sprintf("%d.", i)) + WordStyleBold.Italic(true).Render(p.Point+" "))
		b.WriteString(WordStyle.Render(p.Explanation + "\n"))
		for _, e := range p.SimilarExamples {
			b.WriteString("\t" + DotStyle.Render("→") + MutedStyle.Render("  "+e+"\n"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

type Practice struct {
	PracticeExercises []struct {
		Task   string `json:"task"`
		Answer string `json:"answer"`
	} `json:"practice_exercises"`
}

func (p Practice) Render() string {
	var b strings.Builder

	b.WriteString(WordStyleBold.Render("Practice Exercises:") + "\n")

	if len(p.PracticeExercises) == 0 {
		b.WriteString(MutedStyle.Render("None") + "\n")
	} else {
		for i, ex := range p.PracticeExercises {
			b.WriteString(DotStyle.Render(fmt.Sprintf("%d.", i+1)) + " ")
			b.WriteString(WordStyleBold.Italic(true).Render(ex.Task) + "\n")
			b.WriteString("\t" + MutedStyle.Render("Answer: "+ex.Answer) + "\n\n")
		}
	}

	return b.String()
}

type Usage struct {
	NuanceAndRegister          string   `json:"nuance_and_register"`
	CommonErrors               []string `json:"common_errors"`
	ParaphrasesAndAlternatives []string `json:"paraphrases_and_alternatives"`
	Practice                   Practice
}

func (u Usage) Render() string {
	var b strings.Builder

	b.WriteString(WordStyleBold.Render("Nuances:") + "\n")
	if u.NuanceAndRegister == "" {
		b.WriteString(MutedStyle.Render("None provided") + "\n\n")
	} else {
		for _, line := range strings.Split(u.NuanceAndRegister, "\n") {
			b.WriteString(WordStyle.Render(line) + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(WordStyleBold.Render("Errors:") + "\n")
	if len(u.CommonErrors) == 0 {
		b.WriteString(MutedStyle.Render("None") + "\n\n")
	} else {
		for _, err := range u.CommonErrors {
			b.WriteString(DotStyle.Render("→ ") + WordStyle.Render(err) + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(WordStyleBold.Render("Alternatives:") + "\n")
	if len(u.ParaphrasesAndAlternatives) == 0 {
		b.WriteString(MutedStyle.Render("None") + "\n\n")
	} else {
		for _, p := range u.ParaphrasesAndAlternatives {
			b.WriteString(DotStyle.Render("→ ") + WordStyle.Render(p) + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(u.Practice.Render())
	return b.String()
}
