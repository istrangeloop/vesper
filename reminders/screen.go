// main window for studies
// first of all: choose topic.
// you might have made this choice before, and if so it will be
// suggested here. You can always study something else from your list.

// here you can see and answer cards
// cards that you got wrong will show up to remind you to revise
// them first of all

// do a study session with pomodoros
// you can write notes during a pomodoro
// you can also add other notes about what might be distracting you
// during a pomodoro
// when you finish a session I will show you what you annotated for you to sort it out
// into question cards, or tasks

// see your desk with all your links.

// a screen with your motivations

// Add Finished Thing button - it will go to your achievement wall.

package reminders

import (
	"arthur/utils"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}
type (
	errMsg error
)

func InitialModel() model {
	db := utils.InitDb()
	ta := textarea.New()
	ta.Placeholder = "Send a message to your future self..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	reminders := ShowReminders(db)
	if len(reminders) == 0 {
		vp.SetContent(`Welcome to the notes room!
		Type a message and press Enter to send.`)
	} else {
		vp.SetContent(strings.Join(reminders, "\n"))
	}

	defer db.Close()
	return model{
		textarea:    ta,
		messages:    reminders,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	db := utils.InitDb()
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("Note: ")+m.textarea.Value())
			AddReminder(db, m.textarea.Value(), "support")
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}
	defer db.Close()
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
