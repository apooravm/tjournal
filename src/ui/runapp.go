package ui

import (
	"log"

	api "github.com/apooravm/tjournal/src/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	JournalManage api.JournalDB
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
)

func InitRun(journManage api.JournalDB) error {
	JournalManage = journManage

	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		return err
	} else {
		return nil
	}
}

type JournMessage tea.Msg
