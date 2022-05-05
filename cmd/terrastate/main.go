package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/webhippie/terrastate/pkg/command"
)

func main() {
	if env := os.Getenv("TERRASTATE_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	if err := command.Run(); err != nil {
		os.Exit(1)
	}
}
