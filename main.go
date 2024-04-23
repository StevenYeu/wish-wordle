package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices     []string
	cursor      int
	selected    map[int]struct{}
	secret      string
	secretIndex int
	win         bool
}

func initialModel() model {
	return model{
		choices:     []string{"Grapes", "Oranges", "Berries", "Apples", "Peaches"},
		secretIndex: 3,
		selected:    make(map[int]struct{}),
		secret:      "Apples",
	}
} 

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
			m.cursor = (m.cursor%len(m.choices) + len(m.choices)) % len(m.choices)
		case "down", "j":
			m.cursor++
			m.cursor = m.cursor % len(m.choices)
		case "enter", " ":
			m.selected[m.cursor] = struct{}{}
			if m.choices[m.cursor] == m.secret {
				m.win = true
				return m, tea.Quit
			}

		}
	}
	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#908caa"))
	s := style.Render("Select the secret word")

	s += "\n\n"
	for i, choice := range m.choices {

		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		choiceText := choice
		checked := "[ ]" // not selected
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8"))
		if _, ok := m.selected[i]; ok {
			checked = style.Render("[\u2713]")
			if i != m.secretIndex {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92"))
				checked = style.Render("[x]")
			}
			choiceText = style.Render(choice)
		}
		s += fmt.Sprintf("%s %s %s\n", cursor, checked, choiceText)
	}

	s += "\nPress q to quit.\n"

	if m.win {
		s += "You selected the right word!\n"
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
