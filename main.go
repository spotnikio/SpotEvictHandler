package main

import (
	"os"

	"SpotEvictHandler/internal/pkg/app"
)

func main() {
	// Creatin a new command with parametes
	cmd := app.NewCommand()

	// Execute the command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
