package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.1.3"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GoZilla",
	Long:  `All software has versions. This is GoZilla's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoZilla simple BitTorrent Tracker v%s -- HEAD", VERSION)
	},
}
