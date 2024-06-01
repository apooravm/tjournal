package config

import (
	"fmt"
	"log"

	"github.com/apooravm/tjournal/src/api"
)

func handleCLIArg(cliArg string) {
	switch cliArg {
	case "help":
		fmt.Println("Usage: `tjournal.exe [ARG]` if arg needed\n\nAvailable Args\nhelp   - Display help\ndelete - Delete user config.json")

	case "delete":
		if ConfigFileExists() {
			if err := DeleteConfigFile(); err != nil {
				fmt.Println("Error deleting config")
				return

			} else {
				fmt.Println("Deleted successfully!")
			}
		} else {
			fmt.Println("Config file does not exist")
		}
	}
}

// Print colored error text
func LogColourSprintf(message string, colour string) string {
	switch colour {
	case "red":
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", message)

	case "yellow":
		return fmt.Sprintf("\x1b[33m%s\x1b[0m", message)

	case "green":
		return fmt.Sprintf("\x1b[32m%s\x1b[0m", message)

	case "magenta":
		return fmt.Sprintf("\x1b[35m%s\x1b[0m", message)

	case "cyan":
		return fmt.Sprintf("\x1b[36m%s\x1b[0m", message)

	case "blue":
		return fmt.Sprintf("\x1b[34m%s\x1b[0m", message)

	default:
		return fmt.Sprintf("%s", message)
	}
}

func LogColourPrint(message string, colour string) {
	log.Println(LogColourSprintf(message, colour))
}

func CreateConfigIfNotExist(configName string) {
	if ConfigFileExists() {

	}
}

func ConfigBusiness(configName string, loginEndpoint string) (*LocalConfig, error) {
	if ConfigFileExists() {
		config, err := ReadConfig()
		if err != nil {
			return nil, fmt.Errorf("%s\n", "Error reading config file. "+err.Error())
		}
		return config, nil

	} else {
		email, password := ScanUsernamePassword()
		token, err := api.LoginUser(loginEndpoint, email, password)

		if err != nil {
			// Already know its going to be of type ServerErrorRes thus dont need to check with ok.
			serverErr, _ := err.(api.ServerErrorRes)
			return nil, fmt.Errorf("%s\n", fmt.Sprintf("%d %s %s", serverErr.Code, serverErr.Message, serverErr.Simple))
		}

		if err := CreateConfigFile(token.Token, token.Username); err != nil {
			return nil, fmt.Errorf("%s\n", "Error creating config...")
		}

		config, err := ReadConfig()
		if err != nil {
			return nil, fmt.Errorf("%s\n", "Error readinf config file...")
		}

		return config, nil
	}
}
