//go:build windows

package obs

import (
	"fmt"
	"os"
)

func stopProcessPlatform(proc *os.Process) error {
	if err := proc.Kill(); err != nil {
		return fmt.Errorf("failed to kill OBS: %w", err)
	}
	obsProcess = nil
	return nil
}
