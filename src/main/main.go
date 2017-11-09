package main

import (
	"config"
	"fmt"
	//"os"
	//"os/exec"
	"parser"
	"vm"
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

	parser1.InitVim(runConfig, vm.VM_DEFAULT_TEXT_SIZE, vm.VM_DEFAULT_DATA_SIZE, vm.VM_DEFAULT_STACK_SIZE)

	if !parser1.InitKeywords(runConfig) {
		fmt.Println("ERROR: init keywords failed")
		return
	}

	err = parser1.ReadFile(runConfig.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for !parser1.Eof() {
		parser1.Next(runConfig)
		fmt.Printf("line %d: token = %s", parser1.Line(), parser1.Token().String())
		if parser1.Token() == parser.TOKEN_ID {
			fmt.Printf(", id = %s", parser1.CurrentIdName())
		}
		fmt.Printf("\n")
	}

	fmt.Printf("%.*s", 2, "abc")
}
