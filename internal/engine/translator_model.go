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

func NewTranslatorModel(httpClient *http.Client, k string) *TranslatorModel {
	ta := textarea.New()
	ta.CharLimit = 2000
	ta.Placeholder = "私はバカな男だ"
	ta.Focus()

	//km := ta.KeyMap
	//km.InsertNewline = key.NewBinding(key.WithHelp("ctrl+t", "translate"))
	//ta.KeyMap = km

	return &TranslatorModel{
		ta:   ta,
		pter: 0,
		des:  []domain.TargetLang{domain.TargetJapanese, domain.TargetEnglish, domain.TargetIndonesia},

		sc: domain.NewDeepLClient(k, httpClient),
	}
}

func (im *TranslatorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (im *TranslatorModel) View() string {
	pter := im.pter
	deslen := len(im.des)
	targetLanguage := im.des[pter%deslen].String()
	prev := im.des[(deslen+pter-1)%deslen].String()
	next := im.des[(pter+1)%deslen].String()

	return view.LesterViewStyle.Render(fmt.Sprintf(
		"What do you want to translate to %s?\n\n%s",
		targetLanguage, im.ta.View(),
	)) + view.LesterViewNoteStyle.Render(
		fmt.Sprintf("esc/ctrl+c: exit • ctrl+q: back to menu • shift+left: %s • shift+right: %s • ctrl+t: translate", prev, next),
	)
}

func (im *TranslatorModel) translateCmd(ctx context.Context, lang domain.TargetLang, query string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return switchToLoading{}
		},
		func() tea.Msg {
			res, err := im.sc.Translate(ctx, lang, query)
			if err != nil {
				return switchToError{err}
			}

			return res
		},
	)
}

func (im *TranslatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyShiftLeft:
			im.pter--
			if im.pter < 0 {
				im.pter = len(im.des) - 1
			}

		case tea.KeyShiftRight:
			im.pter++
			if im.pter >= len(im.des) {
				im.pter = 0
			}

		case tea.KeyCtrlT:
			query := im.ta.Value()
			if query == "" {
				return im, nil
			}

			lang := im.des[im.pter%len(im.des)]

			im.ta.Reset()
			return im, im.translateCmd(context.Background(), lang, query)

		case tea.KeyCtrlQ:
			im.ta.Reset()
			return im, func() tea.Msg {
				return switchToMenu{}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return im, tea.Quit

		default:
			if !im.ta.Focused() {
				cmd = im.ta.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	im.ta, cmd = im.ta.Update(msg)
	cmds = append(cmds, cmd)
	return im, tea.Batch(cmds...)
}

func (im *TranslatorModel) Focus() tea.Cmd {
	return im.ta.Focus()
}
