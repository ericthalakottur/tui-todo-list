package main

import (
	"fmt"
	"strings"

	//"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	NUMBER_OF_WINDOWS = 3
	TABLE_WIDTH       = 30
	TABLE_HEIGHT      = 7
)

type newTaskModel struct {
	currentIndex int
	inputs       []textinput.Model
}

type taskListModel struct {
	table table.Model
}

type model struct {
	windowIndex    int
	dbConnection   *DBConnection
	newTaskWindow  newTaskModel
	taskListWindow taskListModel
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func initialModel(dbConnection *DBConnection) model {
	newTaskWindow := newTaskModel{
		inputs: make([]textinput.Model, 2),
	}

	for i := range newTaskWindow.inputs {
		t := textinput.New()
		switch i {
		case 0:
			t.Focus()
			t.Placeholder = "Task"
		case 1:
			t.Placeholder = "Due Date (Format: DD-MM-YYYY)"
		}

		newTaskWindow.inputs[i] = t
	}

	var taskListWindow taskListModel
	taskListWindow.updateTable(dbConnection)

	return model{
		dbConnection:   dbConnection,
		newTaskWindow:  newTaskWindow,
		taskListWindow: taskListWindow,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left":
			m.windowIndex, _ = modulus(m.windowIndex-1, NUMBER_OF_WINDOWS)
			if m.windowIndex == 1 {
				m.taskListWindow.table.Focus()
			} else {
				m.taskListWindow.table.Blur()
			}
		case "right":
			m.windowIndex, _ = modulus(m.windowIndex+1, NUMBER_OF_WINDOWS)
			if m.windowIndex == 1 {
				m.taskListWindow.table.Focus()
			} else {
				m.taskListWindow.table.Blur()
			}
		case "up", "down":
			shift := 1
			if msg.String() == "up" {
				shift = -1
			}
			switch m.windowIndex {
			case 0:
				currentWindow := &m.newTaskWindow
				currentWindow.currentIndex, _ = modulus(currentWindow.currentIndex+shift, len(currentWindow.inputs))
				for i := range currentWindow.inputs {
					currentWindow.inputs[i].Blur()
				}
				currentWindow.inputs[currentWindow.currentIndex].Focus()
			}
		case "enter":
			switch m.windowIndex {
			case 0:
				currentWindow := &m.newTaskWindow
				taskName := currentWindow.inputs[0].Value()
				if taskName != "" {
					completeBy := currentWindow.inputs[1].Value()
					m.dbConnection.CreateNewTask(taskName, completeBy)

					for i := range currentWindow.inputs {
						currentWindow.inputs[i].Reset()
					}
				}
			}
		}

	}

	if m.windowIndex == 0 {
		cmd := m.newTaskWindow.updateInputs(msg)

		return m, cmd
	} else if m.windowIndex == 1 {
		var cmd tea.Cmd

		m.taskListWindow.table, cmd = m.taskListWindow.table.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	var currentWindow strings.Builder

	switch m.windowIndex {
	case 0:
		currentWindow.WriteString("Create New Task\n")
		for i := range m.newTaskWindow.inputs {
			currentWindow.WriteString(m.newTaskWindow.inputs[i].View() + "\n")
		}
	case 1:
		currentWindow.WriteString("Window 2\n")

		m.taskListWindow.updateTable(m.dbConnection)
		currentWindow.WriteString(m.taskListWindow.table.View())
	case 2:
		currentWindow.WriteString("Window 3")
	}
	return currentWindow.String()
}

func (newTaskWindow newTaskModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(newTaskWindow.inputs))

	for i := range newTaskWindow.inputs {
		newTaskWindow.inputs[i], cmds[i] = newTaskWindow.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (taskListWindow *taskListModel) updateTable(dbConnection *DBConnection) {
	columns := []table.Column{
		{Title: "Task #", Width: TABLE_WIDTH},
		{Title: "Task Name", Width: TABLE_WIDTH},
		{Title: "Complete By", Width: TABLE_WIDTH},
	}

	taskList, _ := dbConnection.queries.GetIncompleteTasks(dbConnection.ctx)
	rows := make([]table.Row, len(taskList))
	for i := range len(taskList) {
		task := taskList[i]
		rows[i] = table.Row{fmt.Sprintf("%d", i+1), task.Name, task.CompleteBy.String}
	}

	taskListWindow.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(TABLE_HEIGHT),
	)
}
