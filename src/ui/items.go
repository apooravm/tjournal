package ui

import (
	api "github.com/apooravm/tjournal/src/api"
	"github.com/charmbracelet/bubbles/list"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.desc }

func getItemList(logs *[]api.ReadJournalLogRes) *[]list.Item {
	var items []list.Item
	for _, log := range *logs {
		items = append(items, item{title: log.Title, desc: log.Log + timeStrParser(log.Created_at)})
	}
	return &items
}
