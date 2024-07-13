package exex

import (
	"fmt"
	"os"
	"os/exec"
)

func Exec(command string, path string) {
	cmd := exec.Command(command)
	cmd.Path = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		panic("could not execute command: " + command)
	}
}
