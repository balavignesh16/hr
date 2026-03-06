
package runner

import (
	"os/exec"
	"strconv"
)

func setSysProcAttr(cmd *exec.Cmd) {
}

func killProcess(cmd *exec.Cmd) {
	pid := strconv.Itoa(cmd.Process.Pid)
	killCmd := exec.Command("taskkill", "/T", "/F", "/PID", pid)
	killCmd.Run()
}