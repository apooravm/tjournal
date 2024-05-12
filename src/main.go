package main

import (
	api "github.com/apooravm/tjournal/src/api"
	configMng "github.com/apooravm/tjournal/src/config"
	ui "github.com/apooravm/tjournal/src/ui"
)

var (
	configName    = "tjournalConfig.json"
	URL           = "http://localhost:4000/api/journal/"
	PingUrl       = "http://localhost:4000/api/ping"
	PingUrl2       = "https://multi-serve.onrender.com/api/ping"
	loginEndpoint = "http://localhost:4000/api/user/login"
)

func main() {
	// ConfigBusiness

	// Check if server online
	status, err := api.CheckServerStatus(PingUrl)
	if err != nil {
		configMng.LogColourPrint(err.Error(), "yellow")
		return
	}

	if !status {
		configMng.LogColourPrint("Offline", "red")
		return
	}

	config := configMng.ConfigBusiness(configName, loginEndpoint, PingUrl)
	journalManage := api.JournalDB{Url: URL, Username: config.Username, Token: config.Token}

	if config != nil {
		if err := ui.InitRun(PingUrl, journalManage); err != nil {
			configMng.LogColourPrint(err.Error(), "red")
			return
		}
	}
}
