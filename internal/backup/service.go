package backup

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Ccccraz/cogmoteGO/internal/config"
)

var errBackupRunning = errors.New("backup already running")

type BackupService interface {
	Create(id string, source Source, destination Destination) (Task, error)
	Current() (Task, bool)
}

type service struct {
	mu    sync.RWMutex
	task  *Task
	roots *rootRegistry
}

func newService(sourceRoots, sambaRoots []config.BackupRoot) (*service, error) {
	roots, err := newRootRegistry(sourceRoots, sambaRoots)
	if err != nil {
		return nil, err
	}
	return &service{roots: roots}, nil
}

func (s *service) Create(id string, source Source, destination Destination) (Task, error) {
	sourceRoot, ok := s.roots.source(source.RootID)
	if !ok {
		return Task{}, fmt.Errorf("unknown source root %q", source.RootID)
	}
	if destination.Type != "samba" {
		return Task{}, fmt.Errorf("unsupported destination type %q", destination.Type)
	}
	destinationRoot, ok := s.roots.sambaRoot(destination.RootID)
	if !ok {
		return Task{}, fmt.Errorf("unknown Samba root %q", destination.RootID)
	}
	entries, err := normalizeEntries(source.Entries)
	if err != nil {
		return Task{}, err
	}
	destination.Path, err = normalizeDestinationPath(destination.Path)
	if err != nil {
		return Task{}, err
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	task := Task{
		ID:          id,
		Status:      Running,
		Phase:       PhaseScanning,
		Source:      Source{RootID: source.RootID, Entries: entries},
		Destination: destination,
		Entries:     makePendingEntries(entries),
		CreatedAt:   now,
		StartedAt:   now,
	}
	s.mu.Lock()
	if s.task != nil && s.task.Status == Running {
		activeID := s.task.ID
		s.mu.Unlock()
		return Task{}, fmt.Errorf("%w: task %s is currently running", errBackupRunning, activeID)
	}
	// A new task replaces the last completed task because history is not retained.
	s.task = &task
	s.mu.Unlock()

	go s.execute(id, sourceRoot, destinationRoot)
	return task, nil
}

func makePendingEntries(paths []string) []Entry {
	entries := make([]Entry, 0, len(paths))
	for _, path := range paths {
		entries = append(entries, Entry{Path: path, Status: EntryPending})
	}
	return entries
}

func (s *service) execute(id string, sourceRoot, destinationRoot trustedRoot) {
	defer func() {
		if recovered := recover(); recovered != nil {
			s.finish(id, Failed, fmt.Sprintf("backup panicked: %v", recovered))
		}
	}()

	task, ok := s.Current()
	if !ok || task.ID != id {
		return
	}
	status, err := s.runTransfer(task, sourceRoot, destinationRoot)
	if err != nil {
		s.finish(id, status, err.Error())
		return
	}
	s.finish(id, Succeeded, "")
}

func (s *service) finish(id string, status Status, taskError string) {
	s.update(id, func(task *Task) {
		task.Status = status
		if status == Succeeded {
			task.Phase = PhaseCompleted
		}
		task.Error = taskError
		task.CurrentPath = ""
		task.FinishedAt = time.Now().UTC().Format(time.RFC3339Nano)
	})
}

func (s *service) update(id string, change func(*Task)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.task == nil || s.task.ID != id {
		return
	}
	task := *s.task
	change(&task)
	s.task = &task
}

func (s *service) beginUploading(id string, plan transferPlan) {
	s.update(id, func(task *Task) {
		task.Phase = PhaseUploading
		task.FilesTotal = int64(len(plan.files))
		task.BytesTotal = plan.totalBytes
		for index, entryType := range plan.entryTypes {
			task.Entries[index].Type = entryType
		}
	})
}

func (s *service) setPhase(id string, phase Phase) {
	s.update(id, func(task *Task) {
		task.Phase = phase
		task.CurrentPath = ""
	})
}

func (s *service) setCurrentPath(id, path string) {
	s.update(id, func(task *Task) { task.CurrentPath = path })
}

func (s *service) addTransferredBytes(id string, bytes int64) {
	s.update(id, func(task *Task) { task.BytesTransferred += bytes })
}

func (s *service) completeFile(id string) {
	s.update(id, func(task *Task) {
		task.FilesCompleted++
	})
}

func (s *service) setEntryStatus(id string, index int, status EntryStatus, entryError string) {
	s.update(id, func(task *Task) {
		if index >= len(task.Entries) {
			return
		}
		task.Entries[index].Status = status
		task.Entries[index].Error = entryError
	})
}

func (s *service) Current() (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.task == nil {
		return Task{}, false
	}
	return cloneTask(*s.task), true
}
