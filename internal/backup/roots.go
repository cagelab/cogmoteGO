package backup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ccccraz/cogmoteGO/internal/config"
)

type trustedRoot struct {
	id   string
	path string
}

type rootRegistry struct {
	sources map[string]trustedRoot
	samba   map[string]trustedRoot
}

func newRootRegistry(sourceRoots, sambaRoots []config.BackupRoot) (*rootRegistry, error) {
	sources, err := validateRoots(sourceRoots)
	if err != nil {
		return nil, fmt.Errorf("invalid backup source roots: %w", err)
	}
	samba, err := validateRoots(sambaRoots)
	if err != nil {
		return nil, fmt.Errorf("invalid Samba roots: %w", err)
	}
	return &rootRegistry{sources: sources, samba: samba}, nil
}

func (r *rootRegistry) source(id string) (trustedRoot, bool) {
	root, ok := r.sources[id]
	return root, ok
}

func (r *rootRegistry) sambaRoot(id string) (trustedRoot, bool) {
	root, ok := r.samba[id]
	return root, ok
}

func validateRoots(values []config.BackupRoot) (map[string]trustedRoot, error) {
	roots := make(map[string]trustedRoot, len(values))
	for _, value := range values {
		if !config.IsValidBackupRootID(value.ID) {
			return nil, fmt.Errorf("invalid root id %q", value.ID)
		}
		if _, exists := roots[value.ID]; exists {
			return nil, fmt.Errorf("duplicate root id %q", value.ID)
		}
		if !filepath.IsAbs(value.Path) {
			return nil, fmt.Errorf("root %q path must be absolute", value.ID)
		}
		// os.Root keeps subsequent file operations within the configured root.
		root, err := os.OpenRoot(value.Path)
		if err != nil {
			return nil, fmt.Errorf("open root %q: %w", value.ID, err)
		}
		root.Close()
		roots[value.ID] = trustedRoot{id: value.ID, path: value.Path}
	}
	return roots, nil
}
