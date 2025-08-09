package engine

import (
	"github.com/ziliscite/dictionary-cli/internal/domain"
)

type AppState int

const (
	StateMenu AppState = iota
	StateSearch
	StateLoading
	StateDictionaryList
	StateDetail
	StateTranslate
	StateTranslateDetail
)

type switchToSearch struct{}
type switchToDictionaryNew struct {
	res []domain.Information
}
type switchToDictionaryOld struct{}
type switchToDetail struct {
	res *domain.Information
}
type switchToLoading struct{}
type switchToError struct {
	err error
}
type switchToTranslate struct{}
type switchToTranslateDetail struct {
	res []domain.Translation
}
type switchToMenu struct{}
