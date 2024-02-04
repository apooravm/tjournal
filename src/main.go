package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

func main() {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Could not retrieve exec path")
		return
	}
	pathParts := strings.Split(execPath, "\\")
	pathParts[len(pathParts)-1] = "tjournalConfig.json"
	configJsonPath = strings.Join(pathParts, "\\")

	localConfig = configReadingBusiness()

	// if len(os.Args) == 1 {
	// 	fmt.Println("tjournal.exe <your_log>")
	// 	return
	// }
	// logMessage := strings.Join(os.Args[1:], " ")

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
