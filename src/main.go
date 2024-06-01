package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	api "github.com/apooravm/tjournal/src/api"
	configMng "github.com/apooravm/tjournal/src/config"
	ui "github.com/apooravm/tjournal/src/ui"
)

var (
	configName = "tjournalConfig.json"
	base       = "https://multi-serve.onrender.com"
	PingRoute  = "/api/ping"
	JournRoute = "/api/journal/"
	LoginRoute = "/api/user/login"
	// App states: "quick_save", "quick_view", "tui_view", "tui_save"
	AppState      = ""
	NewLogMessage = ""
	// Some cli args need the main func to return immediately. Toggle this flag for that.
	return_flag = false
)

func handleCLIArg(cliArg []string) {
	switch cliArg[0] {
	case "-help":
		fmt.Println(`Usage: 'tjournal.exe [ARG]' if arg needed
\n
\nAvailable Args
\nnew    - New Log. Usage: 'tjournal.exe -new <YOUR_LOG>. Note: The tags and title are defaulted to 'Quick Note' and 'quick'.
\ndelete - Delete user config.json
\nhelp   - Display help`)
		return_flag = true

	case "-new":
		NewLogMessage = strings.Join(cliArg[1:], " ")
		AppState = "quick_save"

	case "-recent":
		AppState = "quick_view"

	case "-delete":
		if configMng.ConfigFileExists() {
			if err := configMng.DeleteConfigFile(); err != nil {
				configMng.LogColourPrint("\nError deleting config\n", "red")

			} else {
				configMng.LogColourPrint("\nConfig deleted successfully!\n", "green")

			}
		} else {
			configMng.LogColourPrint("\nConfig file does not exist\n", "yellow")

		}
		return_flag = true

	default:
		configMng.LogColourPrint("\nTry tjournal.exe -help\n", "cyan")
	}

}

func main() {
	// ConfigBusiness
	if false {
		base = "http://localhost:4000"
	}

	exePath, err := os.Executable()
	if err != nil {
		configMng.LogColourPrint("Error locating exec file", "yellow")
		return
	}

	exeDir := filepath.Dir(exePath)
	configJsonPath := filepath.Join(exeDir, configName)

	configMng.ConfigPath = configJsonPath

	if len(os.Args) > 1 {
		handleCLIArg(os.Args[1:])
		if return_flag {
			return
		}

	} else {
		AppState = "tui_view"
	}

	// Check internet and server status
	if connStatus := api.UserIsConnected(); !connStatus {
		configMng.LogColourPrint("No internet", "red")
		return
	}

	status, err := api.CheckServerStatus(base + PingRoute)
	if err != nil {
		configMng.LogColourPrint(err.Error(), "yellow")
		return
	}

	if !status {
		configMng.LogColourPrint("Server Offline", "red")
		return
	}

	// If any error, prints it and throws nil
	config, err := configMng.ConfigBusiness(configName, base+LoginRoute)
	if err != nil {
		configMng.LogColourPrint(err.Error(), "red")
		return
	}

	journalManage := api.JournalDB{Url: base + JournRoute, Username: config.Username, Token: config.Token}
	switch AppState {
	case "quick_view":
		fmt.Println("Quick View")
		logs, err := journalManage.ReadJournalLogs()
		if err != nil {
			configMng.LogColourPrint(err.Error(), "red")
			return
		}
		fmt.Println("")
		for _, log := range *logs {
			fmt.Println(log.Title)
			fmt.Println(log.Log)
			fmt.Println(log.Tags)
			fmt.Println("")
		}

	case "quick_save":
		if NewLogMessage == "" {
			fmt.Println("Need log")
		}
		if NewLogMessage != "" {
			journMsg, err := journalManage.CreateJournalLog(NewLogMessage, "Quick Log", &[]string{"quick"})
			if err != nil {
				configMng.LogColourPrint("Error creating log\n", "red")
			}
			if journMsg.Code != 201 {
				configMng.LogColourPrint(journMsg.Message+"\n", "red")
				return

			} else {
				configMng.LogColourPrint("All good pardner 🤠\n", "green")
			}
		}

	case "tui_view":
		if config != nil {
			if err := ui.InitRun(journalManage); err != nil {
				configMng.LogColourPrint(err.Error(), "red")
				return
			}
		}

	case "tui_save":
		fmt.Println("TUI Save")
	}
}
