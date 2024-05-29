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
func LogColourPrint(message string, colour string) {
	switch colour {
	case "red":
		log.Printf("\x1b[31m%s\x1b[0m", message)

	case "yellow":
		log.Printf("\x1b[33m%s\x1b[0m", message)

	case "green":
		log.Printf("\x1b[32m%s\x1b[0m", message)

	case "magenta":
		log.Printf("\x1b[35m%s\x1b[0m", message)

	case "cyan":
		log.Printf("\x1b[36m%s\x1b[0m", message)

	case "blue":
		log.Printf("\x1b[34m%s\x1b[0m", message)

	default:
		log.Printf("%s", message)
	}
}

func CreateConfigIfNotExist(configName string) {
	if ConfigFileExists() {

	}
}

func ConfigBusiness(configName string, loginEndpoint string) *LocalConfig {
	if ConfigFileExists() {
		config, err := ReadConfig()
		if err != nil {
			LogColourPrint("Error reading config file.", err.Error())
			return nil
		}
		return config

	} else {
		email, password := ScanUsernamePassword()
		token, err := api.LoginUser(loginEndpoint, email, password)
		if err != nil {
			serverErr, ok := err.(api.ServerErrorRes)
			if ok {
				LogColourPrint(fmt.Sprintf("%d %s %s", serverErr.Code, serverErr.Message, serverErr.Simple), "red")
			}
			return nil
		}

		if err := CreateConfigFile(token.Token, token.Username); err != nil {
			fmt.Println("Error creating config...")
			return nil
		}

		config, err := ReadConfig()
		if err != nil {
			fmt.Println("Error reading config file...")
			return nil
		}

		return config
	}
}
