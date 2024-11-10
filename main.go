package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dbConnection, err := initializeDB()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(initialModel(dbConnection), tea.WithAltScreen()).Run(); err != nil {
		log.Print("Error starting program: %s", err)
		os.Exit(1)
	}
}
