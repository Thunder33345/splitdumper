package main

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
	"time"
)

var spinnerStyle = lipgloss.NewStyle().Bold(true)

func newModel(w, h int) model {
	m := model{
		width:    w,
		height:   h,
		keyMap:   DefaultKeyMap,
		limit:    3,
		liveSeen: make([]hookMsg, 5),
		liveDest: map[string]struct{}{},
	}

	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = spinnerStyle
	m.spinner = s

	ti := textinput.New()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.Placeholder = "URL"
	ti.Focus()
	wi := 30
	if w < wi {
		wi = w
	}
	ti.Width = wi

	m.input = ti

	sw := stopwatch.NewWithInterval(time.Millisecond * 100)
	m.stopwatch = sw
	return m
}

type model struct {
	width  int
	height int
	//keyMap the default keybindings
	keyMap KeyMap

	spinner   spinner.Model
	input     textinput.Model
	stopwatch stopwatch.Model

	//ch is a channel used to communicate with the dump func
	ch chan tea.Msg
	//cancel is the function to cancel the dump
	cancel func()

	//dumping weather currently it is dumping
	dumping bool
	//done is the dump completed?
	done bool
	//quitting is the app quitting, used for final draw
	quitting bool

	//limit is the times a destination has to be seen
	limit int

	liveSeen []hookMsg
	liveDest map[string]struct{}

	prevUrl   string
	prevLimit int
	urls      []string
	err       error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.HideCursor,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case hookMsg:
		m.liveSeen = append(m.liveSeen[1:], msg)
		m.liveDest[msg.url] = struct{}{}
		return m, waitForActivity(m.ch)
	case resultMsg:
		m.dumping = false
		m.input.Focus()
		m.prevLimit = m.limit
		m.urls = msg.urls
		m.err = msg.err
		m.done = true
		return m, tea.Batch(m.stopwatch.Stop(), waitForActivity(m.ch))
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Quit) {
			if m.cancel != nil {
				m.cancel()
			}
			m.quitting = true
			return m, tea.Quit
		}

		if key.Matches(msg, m.keyMap.Execute) {
			if m.dumping { //stop dumping
				m.dumping = false
				m.input.Focus()
				m.cancel()
				return m, m.stopwatch.Stop()
			} else { //start dumping
				m.done = false
				m.dumping = true
				m.liveSeen = make([]hookMsg, len(m.liveSeen))
				m.liveDest = map[string]struct{}{}
				m.prevLimit = 0
				m.prevUrl = m.input.Value()
				m.urls = []string{}
				m.err = nil

				m.input.Blur()
				m.ch = make(chan tea.Msg, 5)
				m.cancel = dump(m.ch, m.input.Value(), m.limit, time.Millisecond*1200, time.Millisecond*10)
				return m, tea.Batch(m.stopwatch.Reset(), m.stopwatch.Start(), waitForActivity(m.ch))
			}
		}
		if !m.dumping {
			switch {
			case key.Matches(msg, m.keyMap.Incr):
				if m.limit < 10 {
					m.limit++
				}
			case key.Matches(msg, m.keyMap.Decr):
				if m.limit > 1 {
					m.limit--
				}
			}
		}
		if msg.Type == tea.KeyRunes {
			m.input.SetValue("") //clear input for pastes
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg: //resize still breaks if you expand the window size bigger then the initial state
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	}

	var cmdInput, cmdStopwatch tea.Cmd

	//this needs to handle more than input for the blink to work
	m.input, cmdInput = m.input.Update(msg)
	m.stopwatch, cmdStopwatch = m.stopwatch.Update(msg)
	return m, tea.Batch(cmdInput, cmdStopwatch)
}

func (m model) View() string {
	if m.quitting {
		return m.quitView()
	}
	s := bytes.Buffer{}
	s.WriteString("Split Dumper\n")
	s.WriteString("URL:" + m.input.View() + "\n")

	help := ""
	if !m.dumping {
		help = fmt.Sprintf("<%s/%s>", m.keyMap.Decr.Help().Key, m.keyMap.Incr.Help().Key)
	}
	s.WriteString(fmt.Sprintf("Limit: %d %s\n", m.limit, help))

	s.WriteString("Status: ")
	switch {
	case m.dumping:
		s.WriteString(fmt.Sprintf("Dumping %s <%s>", m.spinner.View(), m.keyMap.Execute.Help().Key))
	default:
		s.WriteString(fmt.Sprintf("Idle <%s>", m.keyMap.Execute.Help().Key))
	}
	s.WriteString("\n")

	if m.dumping || m.done {
		s.WriteString(fmt.Sprintf("Elapsed: %s\n", m.stopwatch.View()))
	}
	if m.dumping {
		s.WriteString(fmt.Sprintf("Live (%d urls):\n", len(m.liveDest)))
		for _, msg := range m.liveSeen {
			if msg.seen == 0 || msg.url == "" {
				continue
			}
			s.WriteString(fmt.Sprintf("Saw %d: %s\n", msg.seen, msg.url))
		}
	}
	if m.done {
		s.WriteString(m.resultView())
	}
	return wrap.String(s.String(), m.width)
}

func (m model) quitView() string {
	s := bytes.Buffer{}
	s.WriteString("SplitDumper (Exited)\n")
	s.WriteString(m.resultView())
	return wrap.String(s.String(), m.width)
}
func (m model) resultView() string {
	s := bytes.Buffer{}
	if m.err != nil {
		s.WriteString(fmt.Sprintf("Error: %s\n", m.err))
		s.WriteString("(Partial Results) ")
	}
	if m.prevUrl != "" {
		s.WriteString(fmt.Sprintf("On %s\n", m.prevUrl))
		s.WriteString(fmt.Sprintf("Limit %d Found %d URLs:\n", m.prevLimit, len(m.urls)))
		for _, url := range m.urls {
			s.WriteString(fmt.Sprintf("%s\n", url))
		}
	}
	return s.String()
}
