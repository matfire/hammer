package exec

import (
	"fmt"
	"os/exec"
)

// TODO send writer to log file instead of stdout, stderr
func Exec(command string, path string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
