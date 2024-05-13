package ui

import (
	"strings"

	api "github.com/apooravm/tjournal/src/api"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func GetData() tea.Msg {
	// status, err := api.CheckServerStatus(PingUrl)
	// if err != nil {
	// 	return api.JournError{Code: 400, Message: "Error connecting to the server" + err.Error()}
	// }

	logs, err := JournalManage.ReadJournalLogs()
	if err != nil {
		serverErr, ok := err.(api.ServerErrorRes)
		if !ok {
			return api.JournError{Code: 400, Message: err.Error(), Simple: "Error connecting to the server"}
		}
		return api.JournError{Code: 400, Message: "Error connecting to the server | " + serverErr.Message}
	}

	return logs
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyleTabs      = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type model struct {
	statusCode int
	logs       *[]api.ReadJournalLogRes
	list       list.Model

	filterAgainst string

	inputStyle lipgloss.Style

	quitting bool
	err      error

	tabs         []string
	tabContent   []string
	activeTabIdx int
}

func InitialModel() model {
	var items []list.Item

	del := list.NewDefaultDelegate()
	del.ShowDescription = true
	del.SetHeight(5)
	// del.SetSpacing(1)

	logList := list.New(items, del, 0, 0)

	m := model{list: logList}
	m.list.Title = "Journal Logs"
	m.list.SetSpinner(spinner.Line)

	// m.list.SetShowHelp(false)
	m.filterAgainst = "title"

	m.tabs = []string{"Read Logs", "Create Log"}
	m.tabContent = []string{"", ""}
	m.activeTabIdx = 0

	m.inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7"))
	return m
}

func (m model) Init() tea.Cmd {
	m.tabContent[0] = m.JournalLogReadView()
	return tea.Batch(GetData, func() tea.Msg {
		var msg JournMessage = "startspinner"
		return msg
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "right", "l", "n", "tab":
			m.activeTabIdx = min(m.activeTabIdx+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTabIdx = max(m.activeTabIdx-1, 0)
			return m, nil
		}

	case api.JournError:
		m.err = msg
		return m, nil

	case *[]api.ReadJournalLogRes:
		m.statusCode = 200
		m.logs = msg
		m.list.StopSpinner()

		newKeyBindings := []key.Binding{key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "New log"))}

		m.list.AdditionalShortHelpKeys = func() []key.Binding {
			return newKeyBindings
		}

		return m, m.list.SetItems(*getItemList(m.logs))
		// return m, m.list.StartSpinner()

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		// Helper display
		// m.help.Width = msg.Width

	case JournMessage:
		if msg == "startspinner" {
			return m, m.list.StartSpinner()
		}
	}

	m.tabContent[0] = m.JournalLogReadView()
	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	return m, listCmd
}

func (m model) JournalLogReadView() string {
	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")

	if m.err != nil {
		return m.err.Error()
	}

	if m.quitting {
		return "Bye!\n"
	}

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, m.list.View(), question))
}

func (m model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTabIdx
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.tabContent[m.activeTabIdx]))
	return docStyleTabs.Render(doc.String())

}
