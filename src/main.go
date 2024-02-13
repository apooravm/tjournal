package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	api "github.com/apooravm/tjournal/src/api"
	configMng "github.com/apooravm/tjournal/src/config"

	"github.com/charmbracelet/bubbles/spinner"
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

var docStyle = lipgloss.NewStyle().Margin(1, 4)

var (
	journalManage api.JournalDB
)

type item struct {
	title, description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description + "\n4th july 2011" }
func (i item) FilterValue() string { return i.title }

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
	journalManage = api.JournalDB{Url: URL, Username: config.Username, Password: config.Password}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		return
	}
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return tea.Batch(m.spinner.Tick, GetData)
	}
}

func GetData() tea.Msg {
	// status, err := journalManage.CheckServerStatus()
	// if err != nil {
	// 	fmt.Println("Error checking status", err.Error())
	// 	return JournError{Code: 400, Message: "Error connecting to the server"}
	// }

	// if status.Simple == "bad" {
	// 	return JournError{Code: 500, Message: "Server Offline"}
	// }

	logs, err := journalManage.ReadJournalLogs()
	if err != nil {
		fmt.Println("Something went wrong", err.Error())
		return JournError{Code: 400, Message: "Error connecting to the server"}
	}

	return logs
}

type JournError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e JournError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

type model struct {
	statusCode int
	spinner    spinner.Model
	quitting   bool
	// list       list.Model

	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	// case tea.WindowSizeMsg:
	// 	h, v := docStyle.GetFrameSize()
	// 	m.list.SetSize(msg.Width-h, msg.Height-v)

	case JournError:
		m.err = msg

	case *[]api.ReadJournalLogRes:
		m.statusCode = 200

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	}

	var cmd tea.Cmd
	// m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	s := fmt.Sprintf("\n\n Loading %s\n", m.spinner.View())

	if m.statusCode > 0 {
		s = "got"
	}
	return s
}

// func main_listitem() {
// 	items := []list.Item{
// 		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
// 		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
// 		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
// 		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
// 		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
// 		item{title: "Raspberry Pi’s", description: "I have ’em all over my house"},
// 	}

// 	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
// 	m.list.Title = "My Fave Things"

// 	p := tea.NewProgram(m, tea.WithAltScreen())

// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}
// }
