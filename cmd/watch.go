package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
	"github.com/siegerts/drip/utils"
	"github.com/spf13/cobra"
)

var watchDir string
var entryPoint string
var subDirsToSkip []string
var displayRoutes bool
var hostValue string
var portValue int
var absoluteHost bool
var routeFilter string

func init() {
	watchCmd.Flags().StringVarP(&watchDir, "dir", "d", "", "Source directory to watch")
	watchCmd.Flags().StringVarP(&entryPoint, "entry", "e", "entrypoint.r", "Plumber application entrypoint file")
	watchCmd.Flags().StringSliceVarP(&subDirsToSkip, "skip", "s", []string{".Rproj.user"}, "A comma-separated list of directories to not watch.")
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
		if watchDir != "" {
			if _, err := os.Stat(watchDir); os.IsNotExist(err) {
				fmt.Println("==> Exiting: Directory does not exist")
				os.Exit(1)
			}
			Watch(watchDir)
		} else {
			// watch current
			Watch(".")
		}
	},
}

// Watch is the default explicit run function
func Watch(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		for _, subDir := range subDirsToSkip {
			if info.IsDir() && info.Name() == subDir {
				fmt.Printf("skipping dir without errors: %+v \n", info.Name())
				return filepath.SkipDir
			}
			if info.IsDir() {
				return watcher.Add(path)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	debounced := debounce.New(100 * time.Millisecond)

	var pid int
	plumb := func() {
		// the process needs stopped if it is running
		// if windows, then kill, not interupt
		// if runtime.GOOS == "windows" {
		// 	err = p.Signal(os.Kill)
		// } else {
		// 	//
		// }

		if pid != 0 {
			p, err := os.FindProcess(pid)
			if err != nil {
				fmt.Println(err)
			}
			err = p.Signal(os.Interrupt)
		}

		var plumber string
		if dir != "." {
			plumber = fmt.Sprintf("%s/%s", dir, entryPoint)
		} else {
			plumber = fmt.Sprintf("%s", entryPoint)
		}

		plumbCmd := exec.Command("Rscript", plumber)

		// redirect child output
		plumbCmd.Stdout = os.Stdout
		plumbCmd.Stderr = os.Stderr
		err := plumbCmd.Start()
		pid = plumbCmd.Process.Pid

		if err != nil {
			fmt.Println(err)
		}

		// Execute command
		utils.PrintCommand(plumbCmd)
		// routes
		if displayRoutes {
			RouteStructure(entryPoint, hostValue, portValue, absoluteHost, routeFilter)
		}
		fmt.Println("watching...")

	}

	// watch for used ports
	// https://github.com/jennybc/googlesheets/issues/343#issuecomment-370202906
	// 	lsof -i :8000
	// COMMAND   PID     USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
	// R       38305 siegerts   27u  IPv4 0x79ee09c2f3031013      0t0  TCP localhost:irdmi (LISTEN)

	// initial watch
	fmt.Println("plumbing...")
	plumb()
	if dir != "." {
		log.Println("watching for changes in ", dir)
	} else {
		log.Println("watching for changes in current directory")
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("renamed file: ", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("removed file: ", event.Name)
				}
				fmt.Println("re-plumbing...")
				debounced(plumb)

			case err := <-watcher.Errors:
				fmt.Println("error: ", err)

			case <-done:
				break
			}

		}
	}()

	<-done
}
