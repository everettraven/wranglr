//go:build linux

package linkopener

import "os/exec"

func open(url string) *exec.Cmd {
	// TODO: This could probably be smarter and do some sort of
	// discovery/analysis of what the actual linux distribution
	// is to call tools that are commonly installed by default
	// on those distributions
	return exec.Command("xdg-open", url)
}
