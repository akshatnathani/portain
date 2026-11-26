package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	PID     int
	Name    string
	Command string
}

func Inspect(pid int) (*ProcessInfo, error) {
	if pid <= 0 {
		return nil, fmt.Errorf("invalid PID")
	}

	base := filepath.Join("/proc", strconv.Itoa(pid))

	nameBytes, err := os.ReadFile(filepath.Join(base, "comm"))
	if err != nil {
		return nil, fmt.Errorf("cannot read process: %w", err)
	}

	cmdBytes, err := os.ReadFile(filepath.Join(base, "cmdline"))
	if err != nil {
		return nil, fmt.Errorf("cannot read command: %w", err)
	}

	return &ProcessInfo{
		PID:     pid,
		Name:    strings.TrimSpace(string(nameBytes)),
		Command: strings.ReplaceAll(string(cmdBytes), "\x00", " "),
	}, nil
}
