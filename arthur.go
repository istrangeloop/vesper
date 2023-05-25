package main

import (
	"arthur/reminders"
	"arthur/utils"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
)

func main() {

	parser := argparse.NewParser("arthur", "I help you to do your tasks and study while ADHD.\n arthur r: see and write new reminder")
	//s := parser.String("r", "reminders", &argparse.Options{Required: true, Help: "String to print"})
	var s *string = parser.SelectorPositional([]string{"remind", "study"}, nil)
	var parseerr = parser.Parse(os.Args)

	if parseerr != nil {
		fmt.Print(parser.Usage(parseerr))
	}

	fmt.Println(*s)

	db := utils.InitDb()
	fmt.Println("Connected!")

	// utils.AddReminder(db, "hiiii", "support")
	// utils.AddReminder(db, "howareyou today?", "support")
	var rem = []string{}
	rem = reminders.ShowReminders(db)
	fmt.Println(rem)
	defer db.Close()
	err := db.Ping()
	utils.CheckError(err)

	p := tea.NewProgram(reminders.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oh noes! I errored: %v", err)
		os.Exit(1)
	}
}
