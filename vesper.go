package main

import (
	"vesper/flashcards"
	//"fmt"
	"log"
	"os"

	//"github.com/akamensky/argparse"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")

	p := tea.NewProgram(flashcards.InitialcardModel())
	if _, err = p.Run(); err != nil {
		log.Fatalf("err: %w", err)
		os.Exit(1)
	}
	defer f.Close()
}
