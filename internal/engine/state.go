package engine

import (
	"github.com/ziliscite/dictionary-cli/internal/domain"
)

type AppState int

const (
	StateSearch AppState = iota
	StateLoading
	StateDictionaryList
	StateDetail
	StateTranslate
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
type switchToTranslate struct {
}
