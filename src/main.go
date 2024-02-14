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

type JournMessage tea.Msg

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func timeStrParser(timestr string) string {
	// 2024-02-04T16:17:54.361333+00:00
	y, m, d, time := timestr[0:4], timestr[5:7], timestr[8:10], timestr[11:16]
	return fmt.Sprintf("\n%s, %s-%s-%s", time, d, m, y)
}

func getItemList(logs *[]api.ReadJournalLogRes) *[]list.Item {
	var items []list.Item
	for _, log := range *logs {
		items = append(items, item{title: log.Title + timeStrParser(log.Created_at), desc: log.Log})
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
	list       list.Model

	err error
}

func initialModel() model {
	var items []list.Item

	del := list.NewDefaultDelegate()
	del.SetHeight(4)
	del.SetSpacing(1)

	logList := list.New(items, del, 0, 0)

	m := model{list: logList}
	m.list.Title = "Journal Logs"
	m.list.SetSpinner(spinner.Line)
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(GetData, func() tea.Msg {
		var msg JournMessage = "startspinner"
		return msg
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case api.JournError:
		m.err = msg
		return m, nil

	case *[]api.ReadJournalLogRes:
		m.statusCode = 200
		m.logs = msg
		m.list.StopSpinner()
		return m, m.list.SetItems(*getItemList(m.logs))
		// return m, m.list.StartSpinner()

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case JournMessage:
		if msg == "startspinner" {
			return m, m.list.StartSpinner()
		}
	}

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	return m, listCmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	return docStyle.Render(m.list.View())
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
