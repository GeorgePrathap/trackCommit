package main

import (
	"fmt"
	"os"

	"trackCommit/commands"
	"trackCommit/config"
	"trackCommit/db"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "trackCommit"}

	//connect db
	config.DbConnection = *db.ConnectDB()
	defer db.CloseDatabase(&config.DbConnection)

	rootCmd.AddCommand(getListCommand(), commands.GetAddCommand(), getDeleteCommand(), getMoveCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
