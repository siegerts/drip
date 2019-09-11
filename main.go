package main

import (
	"log"
	"os"

	"github.com/siegerts/drip/cmd"
)

func main() {
	log.SetOutput(os.Stdout)
	cmd.Execute()

}
