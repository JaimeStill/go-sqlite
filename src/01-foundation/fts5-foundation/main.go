package main

import (
	"os"

	"github.com/jaime/go-sqlite/01-foundation/fts5-foundation/commands"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	commands.Root.Init()
	
	err := commands.Root.Command.Execute()
	if err != nil {
		os.Exit(1)
	}
}