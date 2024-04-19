package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Converting to a global module var that can be assigned from configBusiness.go
var (
	ConfigPath string
)

type LocalConfig struct {
	Token    string   `json:"token"`
	Username string   `json:"username"`
	Logs     []string `json:"logs"`
}

func ConfigFileExists() bool {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateConfigFile(token string, username string) error {
	configInit := LocalConfig{
		Token:    token,
		Username: username,
		Logs:     make([]string, 0),
	}

	file, err := os.Create(ConfigPath)
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

func ReadConfig() (*LocalConfig, error) {
	var localConfig LocalConfig

	file, err := os.Open(ConfigPath)
	if err != nil {
		return &localConfig, err
	}

	defer file.Close()

	if err = json.NewDecoder(file).Decode(&localConfig); err != nil {
		return &localConfig, err
	}

	return &localConfig, err
}

func DeleteConfigFile() error {
	if err := os.Remove(ConfigPath); err != nil {
		return err
	}

	return nil
}

// Returns scanned email, password
func ScanUsernamePassword() (string, string) {
	var email string
	var pass string

	fmt.Println("Enter your registered email: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		email = scanner.Text()
	}

	fmt.Println("Enter password: ")
	scanner = bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		pass = scanner.Text()
	}

	return email, pass
}

func LogAppData(appLog string) error {
	byteArr, err := os.ReadFile(ConfigPath)
	if err != nil {
		return errors.New("Error reading file" + err.Error())
	}

	var config LocalConfig
	if err := json.Unmarshal(byteArr, &config); err != nil {
		return fmt.Errorf("error unmarshaling config: %s", err.Error())
	}

	config.Logs = append(config.Logs, appLog)

	updatedByteArr, err := json.Marshal(&config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %s", err.Error())
	}

	if err := os.WriteFile(ConfigPath, updatedByteArr, os.ModePerm); err != nil {
		return fmt.Errorf("error writing to file: %s", err.Error())
	}

	return nil
}
