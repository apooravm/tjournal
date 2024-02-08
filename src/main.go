package main

import (
	"fmt"
	"os"
	"path/filepath"

	api "github.com/apooravm/tjournal/src/api"
	configMng "github.com/apooravm/tjournal/src/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	err            error
	config         *configMng.LocalConfig
	configName     = "tjournalConfig.json"
	configJsonPath string
	URL            = "https://multi-serve.onrender.com/api/journal/"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error locating exec file")
		return
	}
	exeDir := filepath.Dir(exePath)
	configJsonPath = filepath.Join(exeDir, configName)

	if configMng.ConfigFileExists(configJsonPath) {
		config, err = configMng.ReadConfig(configJsonPath)
		if err != nil {
			fmt.Println("Error reading config file...")
			return
		}

	} else {
		username, password := configMng.ScanUsernamePassword()
		if err := configMng.CreateConfigFile(configJsonPath, username, password); err != nil {
			fmt.Println("Error creating config...")
			return
		}

		config, err = configMng.ReadConfig(configJsonPath)
		if err != nil {
			fmt.Println("Error reading config file...")
			return
		}
	}
	journalManage := api.JournalDB{Url: URL, Username: config.Username, Password: config.Password}
	fmt.Println(journalManage)
}

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

func main_listitem() {
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
