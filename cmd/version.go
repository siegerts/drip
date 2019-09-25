package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of drip",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("drip Plumber Application Watcher v0.1.0 -- HEAD")
	},
}
