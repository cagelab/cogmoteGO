package backup

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
)

func normalizeEntries(entries []string) ([]string, error) {
	if len(entries) == 0 {
		return nil, errors.New("at least one source entry is required")
	}
	normalized := make([]string, 0, len(entries))
	seen := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if !filepath.IsLocal(entry) || entry == "." {
			return nil, fmt.Errorf("entry %q must be a non-root relative path", entry)
		}
		entry = filepath.Clean(entry)
		if _, exists := seen[entry]; exists {
			return nil, fmt.Errorf("duplicate entry %q", entry)
		}
		seen[entry] = struct{}{}
		normalized = append(normalized, entry)
	}
	sort.Strings(normalized)
	for i := 1; i < len(normalized); i++ {
		if isDescendant(normalized[i], normalized[i-1]) {
			return nil, fmt.Errorf("overlapping entries %q and %q", normalized[i-1], normalized[i])
		}
	}
	return normalized, nil
}

func normalizeDestinationPath(value string) (string, error) {
	if value == "." {
		return value, nil
	}
	if !filepath.IsLocal(value) {
		return "", fmt.Errorf("destination path %q must be relative", value)
	}
	return filepath.Clean(value), nil
}

func isDescendant(path, parent string) bool {
	relative, err := filepath.Rel(parent, path)
	return err == nil && relative != "." && relative != ".." && !filepath.IsAbs(relative) && !hasParentPrefix(relative)
}

func hasParentPrefix(path string) bool {
	return path == ".." || len(path) > 3 && path[:3] == ".."+string(filepath.Separator)
}
