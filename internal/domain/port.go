package domain

import "context"

type Searcher interface {
	Search(keyword string) ([]Information, error)
}

type Translator interface {
	Translate(ctx context.Context, lang TargetLang, texts ...string) ([]Translation, error)
}
