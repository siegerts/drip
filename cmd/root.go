package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "drip",
	Short: "Watch and reload Plumber applications during development",
	Long: `drip is a utility that will monitor your Plumber applications for any changes in 
your source and automatically restart your server. Perfect for development.`,
	Run: func(cmd *cobra.Command, args []string) {
		// watch current
		cwd, _ := os.Getwd()
		Watch(cwd)
	},
}

// Execute kicks off the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Exiting... ", err)
		os.Exit(1)
	}
}
