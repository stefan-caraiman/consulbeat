package main

import (
	"os"

	"github.com/stefan-caraiman/consulbeat/cmd"

	// Make sure all your modules and metricsets are linked in this file
	_ "github.com/stefan-caraiman/consulbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
