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

func init() {
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch the current directory for any changes",
	Long:  `Watch and rebuild the source if any changes are made across subdirectorys`,
	Run: func(cmd *cobra.Command, args []string) {
		watch()
	},
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	err = filepath.Walk("/Users/siegerts/Documents/projects/feedback2", func(path string, info os.FileInfo, err error) error {
		subDirToSkip := ".Rproj.user"

		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		if info.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})
	if err != nil {
		fmt.Println("ERROR", err)
	}
	log.Println("watching for changes in ", "/Users/siegerts/Documents/projects/feedback2")

	done := make(chan bool)

	debounced := debounce.New(100 * time.Millisecond)

	// this needs moved into seperate flag?
	plumb := func() {
		// requires entrypoint
		fmt.Println("rebuilding...")
		plumber := fmt.Sprintf("/Users/siegerts/Documents/projects/feedback2/entrypoint.r")
		cmd := exec.Command("Rscript", plumber)

		err := cmd.Start()

		// Execute command
		utils.PrintCommand(cmd)

		fmt.Println("watching...")

		if err != nil {
			fmt.Println(err)
		}

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
