package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// PrintCommand will output the current executing command
func PrintCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

// func printError(err error) {
// 	if err != nil {
// 		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
// 	}
// }
