package ggit

import (
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
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
func Branch() string {
	var (
		out string
		err error
	)
	if runtime.GOOS == "windows" {
		if out, err = execShell("git", "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
			return ""
		}
	} else {
		if out, err = execShell("/bin/sh", "-c", "git rev-parse --abbrev-ref HEAD"); err != nil {
			return ""
		}
	}
	return strings.TrimSpace(out)
}

func CommitHash() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}

func CommitShortHash() string {
	if hash := CommitHash(); hash != "" {
		return hash[:7]
	}
	return ""
}

func CommitTimeFormat(layout string) string {
	if ct := VcsTime(); ct != "" {
		if t, e := time.Parse(time.RFC3339, ct); e == nil {
			return t.Format(layout)
		}
	}
	return ""
}

func VcsTime() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}
	return ""
}

func CommitDateTime() string {
	return CommitTimeFormat(time.DateTime)
}

func CommitTime() string {
	return CommitTimeFormat(time.TimeOnly)
}

func CommitDate() string {
	return CommitTimeFormat(time.DateOnly)
}
