package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type LocalConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ConfigFileExists(configpath string) bool {
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateConfigFile(configpath string, username string, password string) error {
	configInit := LocalConfig{
		Username: username,
		Password: password,
	}

	file, err := os.Create(configpath)
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

func ReadConfig(configpath string) (*LocalConfig, error) {
	var localConfig LocalConfig

	file, err := os.Open(configpath)
	if err != nil {
		return &localConfig, err
	}

	defer file.Close()

	if err = json.NewDecoder(file).Decode(&localConfig); err != nil {
		return &localConfig, err
	}

	return &localConfig, err
}

func DeleteConfigFile(configpath string) error {
	if err := os.Remove(configpath); err != nil {
		return err
	}

	return nil
}

// Returns scanned username, password
func ScanUsernamePassword() (string, string) {
	var username string
	var pass string

	fmt.Println("Enter your registered username: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		username = scanner.Text()
	}

	fmt.Println("Enter password: ")
	scanner = bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		pass = scanner.Text()
	}

	return username, pass
}
