package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
)

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

// this needs run in the location of the of the API
func main() {
	// parse args
	// dir
	// assets
	// host
	// port
	// ignore ide files

	// file, err := os.OpenFile("watchr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// log.SetOutput(file)
	log.SetOutput(os.Stdout)

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

	plumb := func() {
		// requires entrypoint
		fmt.Println("rebuilding...")
		plumber := fmt.Sprintf("/Users/siegerts/Documents/projects/feedback2/entrypoint.r")
		cmd := exec.Command("Rscript", plumber)

		err := cmd.Start()

		// Execute command
		printCommand(cmd)

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
