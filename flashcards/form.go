package flashcards

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* FORM MODEL */
type Form struct {
	id       int
	question textinput.Model
	answer   textarea.Model
	keymap   formkeymap
	help     help.Model
}

type formkeymap struct {
	next  key.Binding
	enter key.Binding
	back  key.Binding
}

func NewForm() *Form {
	keymap := formkeymap{
		next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "enter"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
	form := Form{
		// id must be overwritten if form is for editing
		id:       -1,
		question: textinput.New(),
		answer:   textarea.New(),
		keymap:   keymap,
		help:     help.NewModel(),
	}

	form.question.Focus()
	return &form
}

var stylea = lipgloss.NewStyle().
	Width(40).
	Height(10).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")).
	Foreground(lipgloss.Color("227"))

func (m Form) CreateCard() tea.Msg {
	return Card{
		id: m.id, question: m.question.Value(), answer: m.answer.Value(), date: time.Now(),
	}
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.question.Focused() {
				m.question.Blur()
				m.answer.Focus()
				return m, textarea.Blink
			} else {
				models[addCard] = m
				return models[showCard], m.CreateCard
			}
		case "esc":
			return models[showCard], nil
		case "tab":
			if m.question.Focused() {
				m.question.Blur()
				m.answer.Focus()
				return m, textarea.Blink
			} else {
				m.answer.Blur()
				m.question.Focus()
				return m, textinput.Blink
			}
		}

	}
	if m.question.Focused() {
		m.question, cmd = m.question.Update(msg)
		return m, cmd
	} else {
		m.answer, cmd = m.answer.Update(msg)
		return m, cmd
	}
}
func (m Form) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.enter,
		m.keymap.back,
	})
}

func (m Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		stylea.Render(m.question.View()),
		stylea.Render(m.answer.View()),
		m.helpView())
}
