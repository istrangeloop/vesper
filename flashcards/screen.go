package flashcards

import (
	"vesper/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	showCard = 0
	addCard  = 1
)

type cardModel struct {
	deck    Deck
	current int
	show    bool
	keymap  cardkeymap
	help    help.Model
}

type cardkeymap struct {
	add      key.Binding
	edit     key.Binding
	delete   key.Binding
	quit     key.Binding
	conclude key.Binding
	keep     key.Binding
}

/* MODEL MANAGEMENT */
var models []tea.Model

func InitialcardModel() cardModel {
	db := utils.InitDb()
	var cards []Card
	cards = DueCards(db)
	if len(cards) == 0 {
		cards = append(cards, Card{id: 1,
			question: "There are no flashcards to see today!",
			answer:   "add more by pressing a"})
	}
	defer db.Close()
	models = []tea.Model{cardModel{}, NewForm()}
	keymap := cardkeymap{
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		conclude: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "conclude"),
		),
		keep: key.NewBinding(
			key.WithKeys("k", " "),
			key.WithHelp("k", "keep this note"),
		),
	}

	return cardModel{
		deck:    cards,
		current: 0,
		keymap:  keymap,
		help:    help.NewModel(),
		show:    false,
	}
}

func (m cardModel) Init() tea.Cmd {
	return nil
}

func (m *cardModel) NextQuestion() tea.Msg {
	if m.current == len(m.deck)-1 {
		m.current = 0
	} else {
		m.current++
	}
	return nil
}

var style = lipgloss.NewStyle().
	Width(60).
	Height(10).
	PaddingLeft(2).
	PaddingRight(2).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("86")).
	Foreground(lipgloss.Color("202"))

func (m cardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.edit):
			qa := m.deck[m.current]
			editForm := NewForm()
			editForm.question.SetValue(qa.question)
			editForm.answer.SetValue(qa.answer)
			editForm.id = qa.id
			models[showCard] = m // save the state of the current model
			models[addCard] = editForm
			return models[addCard].Update(nil)
		case key.Matches(msg, m.keymap.add):
			models[showCard] = m // save the state of the current model
			models[addCard] = NewForm()
			return models[addCard].Update(nil)
		case key.Matches(msg, m.keymap.conclude, m.keymap.keep):
			m.show = !m.show
			if m.show {
				m.show = true
				return m, nil
			} else {
				if key.Matches(msg, m.keymap.conclude) {
					db := utils.InitDb()
					m.deck[m.current].SendToFuture(db)
					defer db.Close()
				}
				m.NextQuestion()
				m.show = false
				return m, nil
			}
		case key.Matches(msg, m.keymap.delete):
			db := utils.InitDb()
			m.deck[m.current].Delete(db)
			defer db.Close()

			m.deck = m.deck.Remove(m.current)
		}
	// card returned from form
	case Card:
		db := utils.InitDb()
		card := msg
		// if edit, replace existing task in list
		if card.id > 0 {
			card.Update(db)
			m.deck[m.current] = card
		} else {
			card.id = len(m.deck)
			m.deck = append(m.deck, card)
			card.Add(db)
		}
		defer db.Close()
	}

	return m, nil
}

func (m cardModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.add,
		m.keymap.edit,
		m.keymap.delete,
		m.keymap.conclude,
		m.keymap.keep,
		m.keymap.quit,
	})
}

func (m cardModel) View() string {
	s := ""
	if m.show {
		s = m.deck[m.current].answer
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		style.Render(m.deck[m.current].question),
		style.Render(s),
		m.helpView())
}

func (deck Deck) Remove(i int) []Card {
	deck[i] = deck[len(deck)-1]
	return deck[:len(deck)-1]
}
