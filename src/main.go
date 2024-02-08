package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	configJsonPath = "./tjournalConfig.json"
	localConfig    *LocalConfig
	URL            = "https://multi-serve.onrender.com/api/journal/"
)

type LocalConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func configFileExists() bool {
	if _, err := os.Stat(configJsonPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func createConfigFile(username string, password string) error {
	configInit := LocalConfig{
		Username: username,
		Password: password,
	}

	file, err := os.Create(configJsonPath)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(&configInit, "", "    ")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func readConfig() (*LocalConfig, error) {
	var localConfig LocalConfig

	file, err := os.Open(configJsonPath)
	if err != nil {
		return &localConfig, err
	}

	defer file.Close()

	if err = json.NewDecoder(file).Decode(&localConfig); err != nil {
		return &localConfig, err
	}

	return &localConfig, err

}

func configReadingBusiness() *LocalConfig {
	var localConfig *LocalConfig
	if configFileExists() {
		localConfig, err := readConfig()
		if err != nil {
			fmt.Println("Error reading config file...")
			return localConfig
		}
		return localConfig

	} else {
		fmt.Println("No config file found. Creating one...")
		var username string
		var pass string

		fmt.Println("Enter your registered username: ")
		fmt.Scanln(&username)

		fmt.Println("Enter your registered password: ")
		fmt.Scanln(&pass)

		if err := createConfigFile(username, pass); err != nil {
			fmt.Println("Error Creating config file. ", err.Error())
			return localConfig
		}
		fmt.Println("Config file has been created!")

		localConfig, err := readConfig()
		if err != nil {
			fmt.Println("Error reading config file...")
			return localConfig
		}

		return localConfig
	}
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
