package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete the commit from the target branch",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}
}

func getMoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "move",
		Short: "Move the commit from the target branch",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}
}

func getListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all the commit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}
}
