package commands

import "github.com/spf13/cobra"

// CommandGroup represents a hierarchical command structure with child groups and subcommands
type CommandGroup struct {
	Command     *cobra.Command      // The Cobra command this group represents
	ChildGroups []*CommandGroup     // Child command groups (base commands)
	SubCommands []*cobra.Command    // Direct sub-commands
	FlagSetup   func()             // Flag registration function
}

// Init recursively initializes the command group and all its children
func (cg *CommandGroup) Init() {
	// Set up flags for this command group
	if cg.FlagSetup != nil {
		cg.FlagSetup()
	}

	// Initialize and add child groups
	for _, childGroup := range cg.ChildGroups {
		childGroup.Init()
		cg.Command.AddCommand(childGroup.Command)
	}

	// Add direct subcommands
	for _, subCmd := range cg.SubCommands {
		cg.Command.AddCommand(subCmd)
	}
}