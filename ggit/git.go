package ggit

import (
	"fmt"
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

// Repo ...
func Repo() (repo string, err error) {
	var (
		out string
	)
	if out, err = execShell("/bin/sh", "-c", "git remote -v"); err != nil {
		return "", err
	}
	if repo = out[strings.Index(out, ":")+1 : strings.Index(out, ".git")]; repo == "" {
		err = fmt.Errorf("not found, %s", out)
		return "unkown", err
	}
	return repo, nil
}

// Branch ...
func Branch() (branch string, err error) {
	var (
		out string
	)
	if runtime.GOOS == "windows" {
		if out, err = execShell("cmd", "git rev-parse --abbrev-ref HEAD"); err != nil {
			return "", err
		}
	} else {
		if out, err = execShell("/bin/sh", "-c", "git rev-parse --abbrev-ref HEAD"); err != nil {
			return "", err
		}
	}
	return out, err
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
