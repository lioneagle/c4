package main

import (
	//"config"
	"fmt"
	//"os"
	//"os/exec"
	//"parser"
)

func main() {
	/*
		runConfig := config.RunConfig{}
		runConfig.Parse()

		if !runConfig.Check() {
			config.PrintUsage()
			return
		}

		cmd := exec.Command("cmd", "/C", "cls")
		//cmd := exec.Command("cmd")
		//cmd := exec.Command("")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Println("clr failed: ", err.Error())
		}

		fmt.Println("runConfig =", runConfig)
		fmt.Printf("FUN = %v\n", int(parser.FUN))
		fmt.Printf("\x1b[1;40;31m%s\x1b[1;40;32m%s\n", "testPrintColor", "xx1")

		cmd = exec.Command("cmd", "/C", "color")
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			fmt.Println("color failed: ", err.Error())
		}
	*/
	fmt.Printf("%.*s", 2, "abc")
}
