package main

import (
	"fmt"
	"os/exec"
)

func main() {
	pythonexec := "C:/Users/stanley/Anaconda3/python.exe"
	pythonscript := "C:/Users/stanley/Develop/src/_snippets/pysnippets/testscript.py"
	cmd := exec.Command(pythonexec, pythonscript)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}
