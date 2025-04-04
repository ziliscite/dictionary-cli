package internal

import "fmt"

func DisplayJisho(jisho *Jisho) {
	for _, entry := range jisho.Data {
		fmt.Println(entry.Slug)
		DisplayJapaneseTerms(entry.Japanese)
		DisplaySenses(entry.Senses)
		fmt.Println()
	}
}

func DisplayJapaneseTerms(terms []struct {
	Word    string `json:"word,omitempty"`
	Reading string `json:"reading"`
}) {
	for _, term := range terms {
		if term.Word != "" {
			fmt.Println("Word: ", term.Word)
			fmt.Println("Reading: ", term.Reading)
		} else {
			fmt.Println("Word: ", term.Reading)
		}
	}
}

func DisplaySenses(senses []struct {
	EnglishDefinitions []string `json:"english_definitions"`
	PartsOfSpeech      []string `json:"parts_of_speech"`
}) {
	for _, sense := range senses {
		DisplayEnglishDefinitions(sense.EnglishDefinitions)
		DisplayPartsOfSpeech(sense.PartsOfSpeech)
	}
}

func DisplayEnglishDefinitions(definitions []string) {
	for _, definition := range definitions {
		fmt.Println("English definition: ", definition)
	}
}

func DisplayPartsOfSpeech(partsOfSpeech []string) {
	for _, part := range partsOfSpeech {
		fmt.Println("Parts of speech: ", part)
	}
}
