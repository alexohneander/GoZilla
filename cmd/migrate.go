package cmd

import (
	"github.com/alexohneander/GoZilla/internal/database"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the Database of GoZilla",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		database.Migrate()
	},
}
