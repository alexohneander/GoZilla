package cmd

import (
	"github.com/alexohneander/GoZilla/internal/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up the Database of GoZilla",
	Long:  `Clean up, deletes all Peers from the Databse of GoZilla`,
	Run: func(cmd *cobra.Command, args []string) {
		task.ForceCleanPeers()
	},
}
