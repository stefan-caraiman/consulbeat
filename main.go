package main

import (
	"os"

	"github.com/stefan-caraiman/consulbeat/cmd"

	_ "github.com/stefan-caraiman/consulbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
