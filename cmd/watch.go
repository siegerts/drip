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

func init() {
	watchCmd.Flags().StringVarP(&watchDir, "dir", "d", "", "Source directory to watch")
	watchCmd.Flags().StringVarP(&entryPoint, "entry", "e", "entrypoint.R", "Plumber application entrypoint file")
	watchCmd.Flags().StringSliceVarP(&subDirsToSkip, "skip", "s", []string{"node_modules", ".Rproj.user", ".git"}, "A comma-separated list of directories to not watch.")
	watchCmd.Flags().BoolVar(&displayRoutes, "routes", false, "Display route map alongside file watcher")
	watchCmd.PersistentFlags().StringVar(&hostValue, "host", "127.0.0.1", "Display route endpoints with a specific host")
	watchCmd.PersistentFlags().IntVar(&portValue, "port", 8000, "Display route endpoints with a specific port")
	watchCmd.PersistentFlags().BoolVar(&absoluteHost, "showHost", false, "Display absolute route endpoint in output")
	watchCmd.Flags().StringVarP(&routeFilter, "filter", "f", "", "Filter endpoints by substring match")
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
			app.Watch()

		}
	},
}

// Watch is the default explicit run function
func (a *Application) Watch() {

	watcher, err := fsnotify.NewWatcher()
	a.watcher = watcher
	if err != nil {
		fmt.Printf("error: %s \n", err)
	}
	defer a.watcher.Close()

	err = filepath.Walk(a.dir, func(path string, info os.FileInfo, err error) error {
		var skip bool
		for _, subDir := range a.skipDirs {
			skip = info.IsDir() && info.Name() == subDir
			if skip {
				fmt.Printf("[%s] skipping directory: %+v \n", a.path(), info.Name())
				return filepath.SkipDir
			}

		}

		if info.IsDir() {
			return a.watcher.Add(path)
		}
		return nil

	})

	if err != nil {
		fmt.Printf("[%s] error traversing directory... \n", a.path())
	}

	done := make(chan bool)
	debounced := debounce.New(100 * time.Millisecond)

	// initial watch
	fmt.Printf("[%s] plumbing... \n", a.path())
	a.plumb()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("[%s] modified file: %s\n", a.path(), event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Printf("[%s] renamed file: %s\n", a.path(), event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Printf("[%s] removed file: %s\n", a.path(), event.Name)
				}
				fmt.Printf("[%s] plumbing... \n", a.path())
				debounced(a.plumb)

			case err := <-watcher.Errors:
				fmt.Printf("[%s] error: %s\n", a.path(), err)

			case <-done:
				break
			}

		}

	}()

	<-done

}

func (a *Application) plumb() {

	if a.pid != 0 {
		p, err := os.FindProcess(a.pid)
		if err != nil {
			fmt.Println(err) // fatal?
		}

		if runtime.GOOS == "windows" {
			err = p.Signal(os.Kill)
		} else {
			err = p.Signal(os.Interrupt)
		}

	}

	cmd, err := a.determineEntryPoint()
	if err != nil {
		fmt.Println("Exiting... Entrypoint does not exist")
		os.Exit(1)
	}

	plumbCmd := exec.Command("Rscript", cmd)

	// redirect child output
	plumbCmd.Stdout = os.Stdout
	plumbCmd.Stderr = os.Stderr
	err = plumbCmd.Start()

	if err != nil {
		fmt.Println("Exiting... Make sure that R is installed. drip requires Rscript.")
		os.Exit(1)
	}
	a.pid = plumbCmd.Process.Pid

	fmt.Printf("[%s] running: %s \n", a.path(), strings.Join(plumbCmd.Args, " "))

	if displayRoutes {
		fmt.Printf("[%s] routing... \n", a.path())
		a.RouteStructure()
	}
	fmt.Printf("[%s] watching... \n", a.path())

}

func (a *Application) determineEntryPoint() (string, error) {

	if a.dir != "" {
		file := fmt.Sprintf("%s/%s", a.dir, a.entryPoint)
		err := fileExists(file)

		if err != nil {
			return "", err
		}

		return file, nil
	}

	err := fileExists(a.entryPoint)
	if err != nil {
		return "", err
	}

	return a.entryPoint, nil

}

func fileExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}
	return nil
}
