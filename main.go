package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type model struct {
	Board [][]string
	Tries int
	index int
}

func initialModel() model {
	a := make([][]string, 6)
	for i := range a {
		a[i] = make([]string, 5)
		for j := range a[i] {
			a[i][j] = StyleLetter(" ")
		}
	}

	return model{Board: a}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Get(row, col int) string {
	return m.Board[row][col]
}

func (m *model) Delete() {

	if m.index >= 0 {
		m.index--
	}

	if m.index < 0 {
		m.index = 0
	}
}

func StyleLetter(letter string) string {
	var style = lipgloss.NewStyle().
		PaddingLeft(1)

	return style.Render(letter)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currnetTry := m.Tries
	currentIndex := m.index
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strMsg := msg.String(); strMsg {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			currentChar := m.Get(currnetTry, currentIndex)
			if currentChar != StyleLetter(" ") && currentIndex >= 0 {
				m.Board[currnetTry][currentIndex] = StyleLetter(" ")
				return m, nil
			}
			prevIndex := currentIndex - 1
			if prevIndex >= 0 {
				m.Board[currnetTry][prevIndex] = StyleLetter(" ")
				m.Delete()
				return m, nil
			}

			return m, nil
		case "enter":
			if currentIndex == len(m.Board[0])-1 {
				//TODO: Check word
				currnetTry++
				return m, nil
			}

		case regexp.MustCompile(`^[[:alpha:]]$`).FindString(strMsg):
			if currentIndex < len(m.Board[currnetTry]) {
				m.Board[currnetTry][currentIndex] = StyleLetter(strings.ToUpper(strMsg))
				m.index++
			}
			if m.index >= len(m.Board[currnetTry]) {
				m.index--
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) RenderBoard() string {
	t := table.New().Width(35).
		Border(lipgloss.ThickBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(m.Board...).
		StyleFunc(func(row, col int) lipgloss.Style {
            if row == m.Tries+1 {
                
				return lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder()).BorderForeground(lipgloss.Color("63"))
			}

			return lipgloss.NewStyle()
		})

	return t.Render()
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
	s := style.Render("Select the secret word")
	s += "\n"

	s += m.RenderBoard()

	s += "\nPress ctrl+c to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
