// this one will be our pomodoro session assistant!

// we need a timer. When setting it, user will choose how many and
// the intervals: 30 / 5, 25/5, 40/10, whichever works best for her.

// timer will run and display a cute hourglass.

// when timer is up, a sound will play (can cli do this?) or a
// message will blink.

// you can mark sessions as "fruitful" or "unfruitful"
// so if you did little today you will do it again sooner.
// if you did a lot you will focus on other tasks.

package study

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pomoModel struct {
	focus    time.Duration
	rest     time.Duration
	resting  bool
	show     bool
	timer    timer.Model
	keymap   keymap
	help     help.Model
	quitting bool
}

type keymap struct {
	start  key.Binding
	stop   key.Binding
	reset  key.Binding
	quit   key.Binding
	toggle key.Binding
}

type TimesUp tea.Msg

func InitialpomoModel() pomoModel {
	var focus = time.Second * 10
	var rest = time.Second * 5
	m := pomoModel{
		focus:   focus,
		rest:    rest,
		resting: false,
		show:    false,
		timer:   timer.NewWithInterval(focus, time.Millisecond),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
			toggle: key.NewBinding(
				key.WithKeys("t"),
				key.WithHelp("t", "toggle display"),
			),
		},
		help: help.NewModel(),
	}
	m.keymap.start.SetEnabled(false)
	return m
}

func (m pomoModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m pomoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		m.resting = !m.resting
		if m.resting {
			m.timer.Timeout = m.rest
		} else {
			m.timer.Timeout = m.focus
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.toggle):
			m.show = !m.show
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = time.Second * m.focus
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}
func (m pomoModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.toggle,
		m.keymap.quit,
	})
}

func (m pomoModel) View() string {
	s := ""
	if m.timer.Timedout() {
		s = "All done!"
	}
	s += "\n"
	if m.resting {
		s += "resting time...\n"
	}
	if !m.show {
		s += "End of pomodoro in: "
		s += m.timer.View()

	}
	s += m.helpView()
	return lipgloss.JoinVertical(lipgloss.Top, "hourglass", s)
}
