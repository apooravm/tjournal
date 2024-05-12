package main

import (
	api "github.com/apooravm/tjournal/src/api"
	configMng "github.com/apooravm/tjournal/src/config"
	ui "github.com/apooravm/tjournal/src/ui"
)

var (
	configName = "tjournalConfig.json"
	base       = "https://multi-serve.onrender"
	PingRoute  = "/api/ping"
	JournRoute = "/api/journal/"
	LoginRoute = "/api/user/login"
)

func main() {
	// ConfigBusiness
	if true {
		base = "http://localhost:4000"
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
		configMng.LogColourPrint("Offline", "red")
		return
	}

	// If any error, prints it and throws nil
	config := configMng.ConfigBusiness(configName, base+LoginRoute)
	if config != nil {
		journalManage := api.JournalDB{Url: base + JournRoute, Username: config.Username, Token: config.Token}
		if err := ui.InitRun(journalManage); err != nil {
			configMng.LogColourPrint(err.Error(), "red")
			return
		}
	}
}
