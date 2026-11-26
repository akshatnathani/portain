package process

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/akshatnathani/portain/internal/ports"
)

type TerminationResult struct {
	PID        int
	Terminated bool
	PortFreed  bool
}

func Terminate(pid int, port int, timeout time.Duration) (*TerminationResult, error) {
	if pid <= 0 {
		return nil, fmt.Errorf("invalid PID: %d", pid)
	}

	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", port)
	}

	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}

	processes, err := FindByPort(port)
	if err != nil {
		return nil, fmt.Errorf("inspect port: %w", err)
	}

	ownsPort := false
	for _, p := range processes {
		if p.PID == pid {
			ownsPort = true
			break
		}
	}

	if !ownsPort {
		return nil, fmt.Errorf("process %d is not listening on port %d", pid, port)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("find process: %w", err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return nil, fmt.Errorf("send SIGTERM: %w", err)
	}

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		available, err := ports.IsAvailable(port)
		if err != nil {
			return nil, fmt.Errorf("check port availability: %w", err)
		}

		if available {
			return &TerminationResult{
				PID:        pid,
				Terminated: true,
				PortFreed:  true,
			}, nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf(
		"process %d did not release port %d within %s",
		pid,
		port,
		timeout,
	)
}
