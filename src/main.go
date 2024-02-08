package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 4)

type item struct {
	title, description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description + "\n4th july 2011" }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{
		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
