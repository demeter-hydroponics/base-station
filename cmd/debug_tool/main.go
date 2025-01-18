package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"net/http"
)

func main() {
	log.Info("starting server on :12345")
	log.Info("Connect on ws://<ip>:12345/ws")
	http.HandleFunc("/ws", EchoServer)
	//log.Error("Error in server:", "err", http.ListenAndServe(":12345", nil))

	model := NewModel()
	prog := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := prog.Run(); err != nil {
		log.Error("There was an error with TUI", "err", err)
	}
}
