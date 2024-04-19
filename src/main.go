package main

import (
	api "github.com/apooravm/tjournal/src/api"
	configMng "github.com/apooravm/tjournal/src/config"
	ui "github.com/apooravm/tjournal/src/ui"
)

var (
	configName    = "tjournalConfig.json"
	URL           = "http://localhost:4000/api/journal/"
	PingUrl       = "http://localhost:4000/api/cronping"
	loginEndpoint = "http://localhost:4000/api/user/login"
)

func main() {
	config := configMng.ConfigBusiness(configName, loginEndpoint)
	journalManage := api.JournalDB{Url: URL, Username: config.Username, Token: config.Token}

	if err := ui.InitRun(PingUrl, journalManage); err != nil {
		configMng.LogColourPrint(err.Error(), "red")
		return
	}
}
