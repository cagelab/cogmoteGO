package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const copyBufferSize = 1024 * 1024

type plannedFile struct {
	sourcePath string
	entryIndex int
	hash       string
}

type transferPlan struct {
	files       []plannedFile
	directories []string
	entryTypes  []string
	totalBytes  int64
}

func (s *service) runTransfer(task Task, sourceConfig, destinationConfig trustedRoot) (Status, error) {
	sourceRoot, err := os.OpenRoot(sourceConfig.path)
	if err != nil {
		return Failed, fmt.Errorf("open source root: %w", err)
	}
	defer sourceRoot.Close()
	destinationRoot, err := os.OpenRoot(destinationConfig.path)
	if err != nil {
		return Failed, fmt.Errorf("open Samba root: %w", err)
	}
	defer destinationRoot.Close()

	// Scan once to fix totals for a stable client-side progress calculation.
	plan, err := scanTransfer(task, sourceRoot, destinationRoot)
	if err != nil {
		return Failed, err
	}
	s.beginUploading(task.ID, plan)

	if task.Destination.Path != "." {
		if err := destinationRoot.MkdirAll(task.Destination.Path, 0o755); err != nil {
			return Failed, fmt.Errorf("create destination path: %w", err)
		}
	}
	// Write all data to one staging directory so failed uploads are never published.
	staging := filepath.Join(task.Destination.Path, ".partial-"+task.ID)
	if err := destinationRoot.Mkdir(staging, 0o755); err != nil {
		return Failed, fmt.Errorf("create temporary backup directory: %w", err)
	}
	defer destinationRoot.RemoveAll(staging)

	for _, directory := range plan.directories {
		if err := destinationRoot.MkdirAll(filepath.Join(staging, directory), 0o755); err != nil {
			return Failed, fmt.Errorf("create temporary directory: %w", err)
		}
	}

	activeEntry := -1
	for index := range plan.files {
		file := &plan.files[index]
		if activeEntry != file.entryIndex {
			activeEntry = file.entryIndex
			s.setEntryStatus(task.ID, activeEntry, EntryRunning, "")
		}
		s.setCurrentPath(task.ID, file.sourcePath)
		hash, err := copyFile(sourceRoot, destinationRoot, file.sourcePath, filepath.Join(staging, file.sourcePath), func(delta int64) {
			s.addTransferredBytes(task.ID, delta)
		})
		if err != nil {
			s.setEntryStatus(task.ID, file.entryIndex, EntryFailed, err.Error())
			return Failed, fmt.Errorf("copy %s: %w", file.sourcePath, err)
		}
		file.hash = hash
		s.completeFile(task.ID)
	}

	s.setPhase(task.ID, PhaseVerifying)
	// Verify only staging files because callers guarantee that source data is immutable.
	for _, file := range plan.files {
		s.setCurrentPath(task.ID, file.sourcePath)
		actualHash, err := hashFile(destinationRoot, filepath.Join(staging, file.sourcePath))
		if err != nil {
			s.setEntryStatus(task.ID, file.entryIndex, EntryFailed, err.Error())
			return Failed, fmt.Errorf("verify %s: %w", file.sourcePath, err)
		}
		if actualHash != file.hash {
			err := fmt.Errorf("checksum mismatch: %s", file.sourcePath)
			s.setEntryStatus(task.ID, file.entryIndex, EntryFailed, err.Error())
			return Failed, err
		}
	}

	s.setPhase(task.ID, PhasePublishing)
	published := 0
	for index, entry := range task.Entries {
		s.setEntryStatus(task.ID, index, EntryRunning, "")
		target := filepath.Join(task.Destination.Path, entry.Path)
		parent := filepath.Dir(target)
		if parent != "." {
			if err := destinationRoot.MkdirAll(parent, 0o755); err != nil {
				s.setEntryStatus(task.ID, index, EntryFailed, err.Error())
				return statusAfterPublishFailure(published), fmt.Errorf("create destination parent: %w", err)
			}
		}
		// Recheck before publishing in case an external process created the target.
		if _, err := destinationRoot.Lstat(target); err == nil {
			s.setEntryStatus(task.ID, index, EntryFailed, "destination entry appeared while uploading")
			return statusAfterPublishFailure(published), fmt.Errorf("destination entry appeared while uploading: %s", target)
		} else if !errors.Is(err, os.ErrNotExist) {
			s.setEntryStatus(task.ID, index, EntryFailed, err.Error())
			return statusAfterPublishFailure(published), fmt.Errorf("inspect destination entry: %w", err)
		}
		if err := destinationRoot.Rename(filepath.Join(staging, entry.Path), target); err != nil {
			s.setEntryStatus(task.ID, index, EntryFailed, err.Error())
			return statusAfterPublishFailure(published), fmt.Errorf("publish %s: %w", entry.Path, err)
		}
		published++
		s.setEntryStatus(task.ID, index, EntrySucceeded, "")
	}
	return Succeeded, nil
}

func scanTransfer(task Task, sourceRoot, destinationRoot *os.Root) (transferPlan, error) {
	plan := transferPlan{entryTypes: make([]string, len(task.Entries))}
	for index, entry := range task.Entries {
		target := filepath.Join(task.Destination.Path, entry.Path)
		if _, err := destinationRoot.Lstat(target); err == nil {
			return transferPlan{}, fmt.Errorf("destination entry already exists: %s", target)
		} else if !errors.Is(err, os.ErrNotExist) {
			return transferPlan{}, fmt.Errorf("inspect destination entry: %w", err)
		}

		info, err := sourceRoot.Lstat(entry.Path)
		if err != nil {
			return transferPlan{}, fmt.Errorf("inspect source entry %q: %w", entry.Path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 || (!info.Mode().IsRegular() && !info.IsDir()) {
			return transferPlan{}, fmt.Errorf("source entry %q must be a regular file or directory", entry.Path)
		}

		if info.IsDir() {
			plan.entryTypes[index] = "directory"
			if err := scanDirectory(sourceRoot, entry.Path, index, &plan); err != nil {
				return transferPlan{}, err
			}
		} else {
			plan.entryTypes[index] = "file"
			plan.files = append(plan.files, plannedFile{sourcePath: entry.Path, entryIndex: index})
			plan.totalBytes += info.Size()
		}
	}

	return plan, nil
}

func scanDirectory(root *os.Root, sourcePath string, entryIndex int, plan *transferPlan) error {
	return fs.WalkDir(root.FS(), sourcePath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return fmt.Errorf("symbolic links are not supported: %s", path)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.IsDir() {
			plan.directories = append(plan.directories, path)
			return nil
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported file type: %s", path)
		}
		plan.files = append(plan.files, plannedFile{sourcePath: path, entryIndex: entryIndex})
		plan.totalBytes += info.Size()
		return nil
	})
}

func copyFile(sourceRoot, destinationRoot *os.Root, sourcePath, destinationPath string, report func(int64)) (string, error) {
	input, err := sourceRoot.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer input.Close()
	if err := destinationRoot.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return "", err
	}
	output, err := destinationRoot.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	buffer := make([]byte, copyBufferSize)
	for {
		read, readErr := input.Read(buffer)
		if read > 0 {
			written, writeErr := output.Write(buffer[:read])
			if writeErr != nil {
				output.Close()
				return "", writeErr
			}
			if written != read {
				output.Close()
				return "", io.ErrShortWrite
			}
			if _, err := hash.Write(buffer[:written]); err != nil {
				output.Close()
				return "", err
			}
			report(int64(written))
		}
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			output.Close()
			return "", readErr
		}
	}
	if err := output.Sync(); err != nil {
		output.Close()
		return "", err
	}
	if err := output.Close(); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func hashFile(root *os.Root, path string) (string, error) {
	file, err := root.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func statusAfterPublishFailure(published int) Status {
	if published > 0 {
		return PartiallySucceeded
	}
	return Failed
}
