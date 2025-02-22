package main

import (
    "base-station/internal/handlers"
	"github.com/charmbracelet/log"
)

func main() {
	log.Info("Starting up Base Station")
    handlers.Run()
}
