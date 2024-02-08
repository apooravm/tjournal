package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	configJsonPath = "./tjournalConfig.json"
	localConfig    *LocalConfig
	URL            = "https://multi-serve.onrender.com/api/journal/"
)

type LogReqPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Log      string `json:"log"`
}

type LocalConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JournError struct {
	Err    error
	Simple string
}

func (ce *JournError) Error() string {
	return fmt.Sprintf("%v", ce.Simple)
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

type JournalDB struct {
	Url             string
	LocalConfigPath string
	Username        string
	Password        string
}

func (journal *JournalDB) ReadJournalLogs() *[]JournalLogRes {
	payload, err := json.Marshal(UserAuth{
		Username: journal.Username,
		Password: journal.Password,
	})
	if err != nil {
		fmt.Println("Error marshalling")
		return nil
	}

	req, err := http.NewRequest("GET", URL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return nil
	}

	defer res.Body.Close()

	var journalLogs []JournalLogRes

	if err := json.NewDecoder(res.Body).Decode(&journalLogs); err != nil {
		fmt.Println("Error Unmarshallng data")
		return nil
	}

	return &journalLogs
}

func main1() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Could not retrieve exec path")
		return
	}
	pathParts := strings.Split(execPath, "\\")
	pathParts[len(pathParts)-1] = "tjournalConfig.json"
	configJsonPath = strings.Join(pathParts, "\\")

	localConfig = configReadingBusiness()

	// logs := *ReadJournalLogs()
	// for _, log := range logs {
	// 	fmt.Println("Created At:", log.Created_at)
	// 	fmt.Println("Title:", log.Title)
	// 	fmt.Println("Log:", log.Log)
	// 	fmt.Println("Tags:", log.Tags)
	// 	fmt.Println("Log ID:", log.Log_Id)
	// 	fmt.Println("")
	// }

	DeleteJournalLog(46)

	// var newLog JournalLogRes
	// if len(logs) > 0 {
	// 	newLog = logs[0]
	// 	newLog.Title = "Updated Title"
	// 	UpdateJournalLog(&newLog)
	// }

}

func WriteJournalLog() {
	var log string
	var title string
	var tags_str string

	fmt.Println("log: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		log = scanner.Text()
	}

	fmt.Println("title: ")
	scanner = bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		title = scanner.Text()
	}

	fmt.Println("tags (separated with spaces): ")
	scanner = bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		tags_str = scanner.Text()
	}

	tags := strings.Split(tags_str, " ")

	// fmt.Println("LOG:", log)
	// fmt.Println("TITLE:", title)
	// fmt.Println("TAGS:", tags)
	CreateJournalLog(log, title, &tags)
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JournalLogRes struct {
	Created_at string   `json:"created_at"`
	Log        string   `json:"log_message"`
	Title      string   `json:"title"`
	Tags       []string `json:"tags"`
	Log_Id     int      `json:"log_id"`
}

type CreateJournalLogReq struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Log      string   `json:"log"`
	Tags     []string `json:"tags"`
	Title    string   `json:"title"`
}

func CreateJournalLog(log string, title string, tags *[]string) {
	payload, err := json.Marshal(CreateJournalLogReq{
		Username: localConfig.Username,
		Password: localConfig.Password,
		Log:      log,
		Tags:     *tags,
		Title:    title,
	})
	if err != nil {
		fmt.Println("Error creating data payload")
		return
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}

	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Println("Log created!", res.Status)
	} else {
		fmt.Println("Something went wrong.", res.Status)
	}
}

type UpdateLogReq struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Log      string   `json:"log"`
	Tags     []string `json:"tags"`
	Title    string   `json:"title"`
	Log_Id   int      `json:"log_id"`
}

// Create a copy of the original log obj and edit that itself. This becomes the new log
func UpdateJournalLog(prevLog *JournalLogRes) {
	payload, err := json.Marshal(UpdateLogReq{
		Username: localConfig.Username,
		Password: localConfig.Password,
		Log:      prevLog.Log,
		Tags:     prevLog.Tags,
		Title:    prevLog.Title,
		Log_Id:   prevLog.Log_Id,
	})
	if err != nil {
		fmt.Println("Error marshalling payload")
		return
	}

	req, err := http.NewRequest("PUT", URL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}

	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Println("Updated Successfully")
		return
	} else {
		fmt.Println("Something went wrong")
		return
	}
}

type DeleteJournalLogReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Log_Id   int    `json:"log_id"`
}

func DeleteJournalLog(log_id int) {
	payload, err := json.Marshal(DeleteJournalLogReq{
		Username: localConfig.Username,
		Password: localConfig.Password,
		Log_Id:   log_id,
	})
	if err != nil {
		fmt.Println("Error marshalling payload")
		return
	}

	req, err := http.NewRequest("DELETE", URL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}

	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Println("Deleted Successfully")
		return
	} else {
		fmt.Println("Something went wrong", res.Status)
		return
	}
}

// A simple program demonstrating the spinner component from the Bubbles
// component library.

type errMsg error

type model struct {
	spinner  spinner.Model
	quitting bool

	textarea textarea.Model
	err      error
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textarea.New()
	ti.Placeholder = "How was your day?"
	ti.Focus()
	return model{spinner: s, textarea: ti}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Fetching data...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}
