package ui

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/idleberg/hashman/internal/algo"
	"github.com/idleberg/hashman/internal/hasher"
)

type hashDoneMsg struct {
	results []hasher.Result
}

// SpinnerModel is the bubbletea model for the hashing spinner.
type SpinnerModel struct {
	spinner    spinner.Model
	filePath   string
	algorithms []algo.Algorithm
	maxWorkers int
	results    []hasher.Result
	done       bool
}

// Results returns the hash results after the spinner completes.
func (m SpinnerModel) Results() []hasher.Result {
	return m.results
}

// NewSpinnerModel creates a new spinner model for hashing a file.
func NewSpinnerModel(filePath string, algorithms []algo.Algorithm, maxWorkers int) SpinnerModel {
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.BrightCyan)),
	)
	return SpinnerModel{
		spinner:    s,
		filePath:   filePath,
		algorithms: algorithms,
		maxWorkers: maxWorkers,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.doHash)
}

func (m SpinnerModel) doHash() tea.Msg {
	results := hasher.HashFile(m.filePath, m.algorithms, m.maxWorkers)
	return hashDoneMsg{results: results}
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case hashDoneMsg:
		m.results = msg.results
		m.done = true
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SpinnerModel) View() tea.View {
	if m.done {
		return tea.NewView("")
	}
	count := len(m.algorithms)
	noun := "checksums"
	if count == 1 {
		noun = "checksum"
	}
	return tea.NewView(fmt.Sprintf("%s Calculating %d %s for %q\n", m.spinner.View(), count, noun, m.filePath))
}
