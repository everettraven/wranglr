//go:build windows

package linkopener

import "os/exec"

func open(url string) *exec.Cmd {
	return exec.Command("start", url)
}
