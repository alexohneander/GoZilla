package cmd

import (
	"fmt"
	"os"

	"github.com/alexohneander/GoZilla/internal/http"
	"github.com/alexohneander/GoZilla/internal/task"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gozilla",
	Short: "GoZilla is a very fast and simple BitTorrent Tracker",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		go task.CleanPeers()
		http.Server()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
