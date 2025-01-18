package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"io"
	"net/http"
)

var (
	//RxMessages = make(chan map[string]string, 100)
	TxMessages = make(chan map[string]string, 100)
	model      = NewModel()
	prog       = tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
)

func main() {
	log.SetOutput(io.Discard)
	// start up the server
	go func() {
		//log.Info("starting server on :12345")
		//log.Info("Connect on ws://<ip>:12345/ws")
		http.HandleFunc("/ws", MockServer)
		log.Error("Error in server:", "err", http.ListenAndServe(":12345", nil))
	}()

	if _, err := prog.Run(); err != nil {
		log.Error("There was an error with TUI", "err", err)
	}
}
