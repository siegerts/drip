package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "drip",
	Short: "Watch and automatically restart Plumber applications during development",
	Long: `drip is a utility that will monitor your Plumber applications for any changes in 
your source and automatically restart your server. Perfect for development.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := Application{
			dir:           watchDir,
			entryPoint:    entryPoint,
			skipDirs:      subDirsToSkip,
			displayRoutes: displayRoutes,
			host:          hostValue,
			port:          portValue,
			absoluteHost:  absoluteHost,
			routeFilter:   routeFilter,
			// tunnelPort:    tunnelPort,
			pid: 0,
		}
		// watch current

		cwd, _ := os.Getwd()
		app.dir = cwd
		app.Watch()
	},
}

// Execute kicks off the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Exiting... ", err)
		os.Exit(1)
	}
}
