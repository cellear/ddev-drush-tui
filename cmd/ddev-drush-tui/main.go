package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/cellear/ddev-drush-tui/internal/ddev"
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

	fmt.Printf("Project: %s\n", context.ProjectName)
}
