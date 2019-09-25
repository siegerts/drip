package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// Application is the definition of a Plumber app
// for the purposes of development and testing
type Application struct {
	dir           string
	entryPoint    string
	skipDirs      []string
	displayRoutes bool
	host          string
	port          int
	absoluteHost  bool
	routeFilter   string
	pid           int
	watcher       *fsnotify.Watcher
}

func (app *Application) path() string {
	return filepath.Base(app.dir)
}

var rootCmd = &cobra.Command{
	Use:   "drip",
	Short: "Watch and automatically restart Plumber applications during development",
	Long: `drip is a utility that will monitor your Plumber applications for any changes in 
your source and automatically restart your server.`,
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
			pid:           0,
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
