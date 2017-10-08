package main

import (
	"config"
	"fmt"
	"parser"
)

func main() {
	runConfig := config.RunConfig{}
	runConfig.Parse()

	if !runConfig.Check() {
		config.PrintUsage()
		return
	}

	fmt.Println("runConfig =", runConfig)
	fmt.Printf("FUN = %v\n", int(parser.FUN))
}
