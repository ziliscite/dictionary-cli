package domain

type Searcher interface {
	Search(keyword string) ([]Information, error)
}
