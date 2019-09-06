package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "drip",
	Short: "Watch and rebuild Plumber applications",
	Long: `drip is a utility that will monitor your Plumber applications for any changes in 
your source and automatically restart your server. Perfect for development.
Complete documentation is available at x`,
	Run: func(cmd *cobra.Command, args []string) {
		Watch()
	},
}

// Execute kicks off the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
