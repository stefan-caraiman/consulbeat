package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}