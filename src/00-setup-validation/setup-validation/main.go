package main

import (
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/commands"
	"github.com/jaime/go-sqlite/00-setup-validation/setup-validation/errors"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	commands.Root.Init()
	
	err := commands.Root.Command.Execute()
	if err != nil {
		errors.DisplayError(err)
	}
}
