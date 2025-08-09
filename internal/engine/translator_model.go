package engine

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ziliscite/dictionary-cli/internal/domain"
	"github.com/ziliscite/dictionary-cli/internal/view"
	"net/http"
)

type TranslatorModel struct {
	ta   textarea.Model
	pter int
	des  []domain.TargetLang

	sc domain.Translator
}

func NewTranslatorModel(httpClient *http.Client, key string) *TranslatorModel {
	ta := textarea.New()
	ta.CharLimit = 2000
	ta.Placeholder = ""
	ta.Focus()

	return &TranslatorModel{
		ta:   ta,
		pter: 0,
		des:  []domain.TargetLang{domain.TargetJapanese, domain.TargetEnglish, domain.TargetIndonesia},

		sc: domain.NewDeepLClient(key, httpClient),
	}
}

func (im *TranslatorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (im *TranslatorModel) View() string {
	pter := im.pter
	deslen := len(im.des)
	targetLanguage := im.des[pter%deslen].String()
	prev := im.des[(pter-1)%deslen].String()
	next := im.des[(pter+1)%deslen].String()

	return view.LesterViewStyle.Render(fmt.Sprintf(
		"What do you want to translate to %s?\n\n%s",
		im.ta.View(), targetLanguage,
	)) + view.LesterViewNoteStyle.Render(
		fmt.Sprintf("esc/ctrl+c: exit • <-: %s • ->: %s • enter: translate", prev, next),
	)
}

func (im *TranslatorModel) translateCmd(ctx context.Context, query string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return switchToLoading{}
		},
		func() tea.Msg {
			res, err := im.sc.Translate(ctx, domain.TargetJapanese, query)
			if err != nil {
				return switchToError{err}
			}

			return res
		},
	)
}

func (im *TranslatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			query := im.ta.Value()
			if query == "" {
				return im, nil
			}

			im.ta.Reset()
			return im, im.translateCmd(context.Background(), query)

		case tea.KeyCtrlC, tea.KeyEsc:
			return im, tea.Quit

		case tea.KeyLeft:
			im.pter--
			if im.pter < 0 {
				im.pter = len(im.des) - 1
			}

		case tea.KeyRight:
			im.pter++
			if im.pter >= len(im.des) {
				im.pter = 0
			}

		default:
			panic("unhandled default case")
		}
	}

	im.ta, cmd = im.ta.Update(msg)
	return im, cmd
}

func (im *TranslatorModel) Focus() tea.Cmd {
	return im.ta.Focus()
}
