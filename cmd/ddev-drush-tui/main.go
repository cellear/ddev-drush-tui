package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/cellear/ddev-drush-tui/internal/ddev"
	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/cellear/ddev-drush-tui/internal/tui"
)

func main() {
	context, err := ddev.Detect()
	if err != nil {
		switch {
		case errors.Is(err, ddev.ErrNotInProject):
			fmt.Println("Error: not in a DDEV project directory")
		case errors.Is(err, ddev.ErrDrushUnavailable):
			fmt.Println("Error: drush is not available in this DDEV project")
		default:
			fmt.Printf("Error: %v\n", err)
		}
		os.Exit(1)
	}

	commands, err := drush.ListCommands()
	if err != nil {
		fmt.Printf("Error listing commands: %v\n", err)
		os.Exit(1)
	}

	app := tui.NewApp(context, commands)
	if err := app.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
