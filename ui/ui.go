package ui

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	err                error
	spinner            spinner.Model
	quitting           bool
	LocalPingErr       chan error
	RemotePingErr      chan error
	LocalPingErrCount  int
	RemotePingErrCount int
}

func initialModel(local, remote chan error) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, LocalPingErr: local, RemotePingErr: remote}
}

func (m *model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *model) View() (str string) {
	if m.err != nil {
		return m.err.Error()
	}

	select {
	case err := <-m.LocalPingErr:
		str += fmt.Sprintf("Local Error Found %s", err.Error())
		m.LocalPingErrCount++
	case err := <-m.RemotePingErr:
		str += fmt.Sprintf("Remote Error Found %s", err.Error())
		m.RemotePingErrCount++
	default:
		str = ""
	}
	str += fmt.Sprintf("\n %s Remote Pinging and Waiting for error\n    Local Error:%d Remote Error:%d\n", m.spinner.View(), m.LocalPingErrCount, m.RemotePingErrCount)
	return str
}

func CreateUI(local, remote chan error) error {
	model := initialModel(local, remote)
	p := tea.NewProgram(&model)
	_, err := p.Run()
	return err
}
