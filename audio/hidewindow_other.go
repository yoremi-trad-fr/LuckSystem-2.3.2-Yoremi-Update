//go:build !windows

package audio

import "os/exec"

func hideWindow(cmd *exec.Cmd) {}
