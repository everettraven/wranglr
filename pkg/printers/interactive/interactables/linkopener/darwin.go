//go:build darwin

package linkopener

import "os/exec"

func open(url string) *exec.Cmd {
	return exec.Command("open", url)
}
