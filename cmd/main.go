package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jaredhaight/cryptoclicker/internal"
	"os"
)

func main() {
	app := internal.NewApplication()
	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
