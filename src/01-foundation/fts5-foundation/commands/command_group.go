package commands

import "github.com/spf13/cobra"

// CommandGroup represents a hierarchical command structure
type CommandGroup struct {
	// The Cobra command this group represents
	Command *cobra.Command
	
	// Child command groups (base commands like "document", "user", etc.)
	ChildGroups []*CommandGroup
	
	// Direct sub-commands attached to this group
	SubCommands []*cobra.Command
	
	// Flag registration function for this group and its commands
	FlagSetup func()
}


// Init initializes the command group by registering all commands and flags
func (cg *CommandGroup) Init() {
	cg.RegisterCommands()
	cg.RegisterFlags()
}

// RegisterFlags recursively calls flag setup for this group and all children
func (cg *CommandGroup) RegisterFlags() {
	if cg.FlagSetup != nil {
		cg.FlagSetup()
	}
	
	for _, childGroup := range cg.ChildGroups {
		childGroup.RegisterFlags()
	}
}

// RegisterCommands recursively registers all commands in the hierarchy
func (cg *CommandGroup) RegisterCommands() {
	// Add direct sub-commands to this group's command
	for _, subCmd := range cg.SubCommands {
		cg.Command.AddCommand(subCmd)
	}
	
	// Add child groups to this group's command and recursively register their commands
	for _, childGroup := range cg.ChildGroups {
		cg.Command.AddCommand(childGroup.Command)
		childGroup.RegisterCommands()
	}
}