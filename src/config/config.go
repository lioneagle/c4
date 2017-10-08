package config

import (
	"flag"
	"fmt"
	"os"
)

type RunConfig struct {
	Filename    string
	PrintSource bool
	Debug       bool
}

func (runConfig *RunConfig) Parse() {
	flag.StringVar(&runConfig.Filename, "file", "", "source file")
	flag.BoolVar(&runConfig.PrintSource, "source", false, "print source and assembly")
	flag.BoolVar(&runConfig.Debug, "debug", false, "print executed instructions")

	flag.Parse()
}

func (runConfig *RunConfig) Check() bool {
	_, err := os.Stat(runConfig.Filename)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		fmt.Printf("ERROR: file \"%s\" is not exist\n", runConfig.Filename)
	} else {
		fmt.Printf("ERROR: file \"%s\" is invalid\n", runConfig.Filename)
	}
	return false
}

func PrintUsage() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
