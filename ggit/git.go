package ggit

import (
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
)

// execShell ...
func execShell(command string, arg ...string) (out string, err error) {
	var Stdout []byte
	cmd := exec.Command(command, arg...)
	Stdout, err = cmd.CombinedOutput()
	out = string(Stdout)
	return
}

// Branch ...
func Branch() (branch string, err error) {
	var (
		out string
	)
	if runtime.GOOS == "windows" {
		if out, err = execShell("git", "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
			return "", err
		}
	} else {
		if out, err = execShell("/bin/sh", "-c", "git rev-parse --abbrev-ref HEAD"); err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(out), err
}

func Version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}

func LastCommitTime() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}
	return ""
}
