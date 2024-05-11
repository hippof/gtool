package ggit

import (
	"os/exec"
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
	if out, err = execShell("git", "rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return ""
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
	var (
		out string
		err error
	)
	if out, err = execShell("git", "rev-parse", "HEAD"); err != nil {
		return ""
	}
	return strings.TrimSpace(out)
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
	var (
		out string
		err error
	)
	if out, err = execShell("git", "log", "-1", "--date=format:%Y-%m-%d %H:%M:%S", "--format=%cd"); err != nil {
		return ""
	}
	out = strings.TrimSpace(out)
	if t, e := time.Parse(time.DateTime, out); e == nil {
		return t.Format(time.RFC3339)
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
