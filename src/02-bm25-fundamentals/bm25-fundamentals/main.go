package main

import (
	"os"

	"github.com/jaime/go-sqlite/02-bm25-fundamentals/bm25-fundamentals/commands"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize command structure first
	commands.Root.Init()

	err := commands.Root.Command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
