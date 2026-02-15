//go:build !windows

package obs

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func stopProcessPlatform(proc *os.Process) error {
	pgid, err := syscall.Getpgid(proc.Pid)
	if err != nil {
		return fmt.Errorf("failed to get process group: %w", err)
	}

	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM to OBS process group: %w", err)
	}

	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
				return fmt.Errorf("failed to kill OBS process group after timeout: %w", err)
			}
			obsProcess = nil
			return fmt.Errorf("OBS did not exit gracefully, force killed")
		case <-ticker.C:
			if err := proc.Signal(syscall.Signal(0)); err != nil {
				obsProcess = nil
				return nil
			}
		}
	}
}
