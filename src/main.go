package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	api "github.com/apooravm/tjournal/src/api"
	configMng "github.com/apooravm/tjournal/src/config"

	"github.com/charmbracelet/bubbles/list"
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

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var (
	journalManage api.JournalDB
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func getItemList(logs *[]api.ReadJournalLogRes) *[]list.Item {
	var items []list.Item
	for _, log := range *logs {
		items = append(items, item{title: log.Title, desc: log.Log})
	}
	return &items
}

func GetData() tea.Msg {
	status, err := journalManage.CheckServerStatus()
	if err != nil {
		return api.JournError{Code: 400, Message: "Error connecting to the server" + err.Error()}
	}

	if status.Simple == "bad" {
		return api.JournError{Code: 500, Message: "Server Offline"}
	}

	logs, err := journalManage.ReadJournalLogs()
	if err != nil {
		return api.JournError{Code: 400, Message: "Error connecting to the server" + err.Error()}
	}

	return logs
}

type model struct {
	statusCode int
	logs       *[]api.ReadJournalLogRes

	spinner  spinner.Model
	quitting bool

	list list.Model

	err error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, GetData)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		} else {
			return m, nil
		}
	// case tea.WindowSizeMsg:
	// 	h, v := docStyle.GetFrameSize()
	// 	m.list.SetSize(msg.Width-h, msg.Height-v)

	case api.JournError:
		m.err = msg
		return m, nil

	case *[]api.ReadJournalLogRes:
		m.statusCode = 200
		m.logs = msg
		m.list = list.New(*getItemList(m.logs), list.NewDefaultDelegate(), 0, 0)
		m.list.Title = "Journal Logs"

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		if m.statusCode > 0 {
			m.list.SetSize(msg.Width-h, msg.Height-v)
		}
	}

	if m.statusCode > 0 {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	} else {
		// if data not fetched keep spinning
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	s := fmt.Sprintf("\n\n Loading %s\n\n", m.spinner.View())

	if m.statusCode > 0 {
		// s = fmt.Sprintf("\n\n %v \n\n", m.logs)
		s = docStyle.Render(m.list.View())
	}
	return s
}

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
