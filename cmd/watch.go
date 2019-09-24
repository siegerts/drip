package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	watchDir      string
	entryPoint    string
	subDirsToSkip []string
	displayRoutes bool
	hostValue     string
	portValue     int
	absoluteHost  bool
	routeFilter   string
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

func init() {
	watchCmd.Flags().StringVarP(&watchDir, "dir", "d", "", "Source directory to watch")
	watchCmd.Flags().StringVarP(&entryPoint, "entry", "e", "entrypoint.r", "Plumber application entrypoint file")
	watchCmd.Flags().StringSliceVarP(&subDirsToSkip, "skip", "s", []string{"node_modules", ".Rproj.user"}, "A comma-separated list of directories to not watch.")
	watchCmd.Flags().BoolVar(&displayRoutes, "routes", false, "Display route map alongside file watcher")
	watchCmd.PersistentFlags().StringVar(&hostValue, "host", "127.0.0.1", "Display route endpoints with a specific host")
	watchCmd.PersistentFlags().IntVar(&portValue, "port", 8000, "Display route endpoints with a specific port")
	watchCmd.PersistentFlags().BoolVar(&absoluteHost, "showHost", false, "Display absolute route endpoint in output")
	watchCmd.Flags().StringVarP(&routeFilter, "filter", "f", "", "Filter endpoints by prefix match")
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch the current directory for any changes",
	Long:  `Watch and rebuild the source if any changes are made across subdirectories`,
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

		if app.dir != "" {
			if _, err := os.Stat(app.dir); os.IsNotExist(err) {
				fmt.Println("Exiting... Directory does not exist")
				os.Exit(1)
			}
			app.Watch()
		} else {
			// watch current
			cwd, _ := os.Getwd()
			app.dir = cwd
			if _, err := os.Stat(app.dir); os.IsNotExist(err) {
				fmt.Println("Exiting... Error accessing current directory")
				os.Exit(1)
			}
			app.Watch()

		}
	},
}

// Watch is the default explicit run function
func (app *Application) Watch() {

	watcher, err := fsnotify.NewWatcher()
	app.watcher = watcher
	if err != nil {
		fmt.Printf("error: %s \n", err)
	}
	defer app.watcher.Close()

	err = filepath.Walk(app.dir, func(path string, info os.FileInfo, err error) error {
		var skip bool
		for _, subDir := range app.skipDirs {
			skip = info.IsDir() && info.Name() == subDir
			if skip {
				fmt.Printf("[%s] skipping directory: %+v \n", app.path(), info.Name())
				return filepath.SkipDir
			}

		}

		if info.IsDir() {
			return app.watcher.Add(path)
		}
		return nil

	})

	if err != nil {
		fmt.Printf("[%s] error traversing directory... \n", app.path())
	}

	done := make(chan bool)

	debounced := debounce.New(100 * time.Millisecond)

	plumb := func() {

		if app.pid != 0 {
			p, err := os.FindProcess(app.pid)
			if err != nil {
				fmt.Println(err)
			}
			if runtime.GOOS == "windows" {
				err = p.Signal(os.Kill)
			} else {
				err = p.Signal(os.Interrupt)
			}

		}

		var plumber string
		// refactor this into exists function
		if app.dir != "" {
			if _, err := os.Stat(fmt.Sprintf("%s/%s", app.dir, app.entryPoint)); os.IsNotExist(err) {
				fmt.Println("Exiting... Entrypoint does not exist")
				os.Exit(1)
			}
			plumber = fmt.Sprintf("%s/%s", app.dir, app.entryPoint)
		} else {
			if _, err := os.Stat(app.entryPoint); os.IsNotExist(err) {
				fmt.Println("Exiting... Entrypoint does not exist")
				os.Exit(1)
			}
			plumber = fmt.Sprintf("%s", app.entryPoint)
		}

		plumbCmd := exec.Command("Rscript", plumber)

		// redirect child output
		plumbCmd.Stdout = os.Stdout
		plumbCmd.Stderr = os.Stderr
		err := plumbCmd.Start()
		app.pid = plumbCmd.Process.Pid

		if err != nil {
			fmt.Println("Exiting... Error catching process id")
			os.Exit(1)
		}

		fmt.Printf("[%s] running: %s \n", app.path(), strings.Join(plumbCmd.Args, " "))

		if displayRoutes {
			fmt.Printf("[%s] routing... \n", app.path())
			app.RouteStructure()
		}
		fmt.Printf("[%s] watching... \n", app.path())

	}

	// initial watch
	fmt.Printf("[%s] plumbing... \n", app.path())
	plumb()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("[%s] modified file: %s\n", app.path(), event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Printf("[%s] renamed file: %s\n", app.path(), event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Printf("[%s] removed file: %s\n", app.path(), event.Name)
				}
				fmt.Printf("[%s] plumbing... \n", app.path())
				debounced(plumb)

			case err := <-watcher.Errors:
				fmt.Printf("[%s] error: %s\n", app.path(), err)

			case <-done:
				fmt.Printf("done.\n")
				break
			}

		}
	}()

	<-done
}
