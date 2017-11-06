package main

import (
	"config"
	"fmt"
	//"os"
	//"os/exec"
	"parser"
)

func main() {

	//*
	runConfig := &config.RunConfig{}
	runConfig.Parse()

	if !runConfig.Check() {
		config.PrintUsage()
		return
	}
	var err error

	//cmd := exec.Command("cmd", "/C", "cls")
	//cmd := exec.Command("cmd")
	//cmd := exec.Command("")
	///cmd.Stdout = os.Stdout
	//err = cmd.Run()
	//if err != nil {
	//fmt.Println("clr failed: ", err.Error())
	//}

	fmt.Println("runConfig =", runConfig)
	//fmt.Printf("FUN = %v\n", int(parser.FUN))
	//fmt.Printf("\x1b[1;40;31m%s\x1b[1;40;32m%s\n", "testPrintColor", "xx1")

	//cmd = exec.Command("cmd", "/C", "color")
	//cmd.Stdout = os.Stdout
	//err = cmd.Run()
	//if err != nil {
	//fmt.Println("color failed: ", err.Error())
	//}
	//*/

	parser1 := parser.NewParser()

	err = parser1.ReadFile(runConfig.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for !parser1.Eof() {
		parser1.Next(runConfig)
		fmt.Printf("line %d: token = %s", parser1.Line(), parser1.Token().String())
		if parser1.Token() == parser.ID {
			fmt.Printf(", id = %s", parser1.IdName())
		}
		fmt.Printf("\n")
	}

	fmt.Printf("%.*s", 2, "abc")
}
