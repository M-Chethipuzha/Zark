package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Index represents the staging area.
type Index struct {
	Entries []IndexEntry `json:"entries"`
}

type IndexEntry struct {
	Path     string    `json:"path"`
	Hash     string    `json:"hash"`
	Mode     string    `json:"mode"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

func NewIndex() *Index {
	return &Index{
		Entries: make([]IndexEntry, 0),
	}
}

// Add adds or updates an entry in the index.
func (i *Index) Add(path, hash, mode string, size int64, modified time.Time) {
	// Remove existing entry if present to ensure no duplicates.
	for j, entry := range i.Entries {
		if entry.Path == path {
			i.Entries = append(i.Entries[:j], i.Entries[j+1:]...)
			break
		}
	}

	i.Entries = append(i.Entries, IndexEntry{
		Path:     path,
		Hash:     hash,
		Mode:     mode,
		Size:     size,
		Modified: modified,
	})
}

// Save writes the index to a file.
func (i *Index) Save(path string) error {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	return nil
}

// LoadIndex reads and unmarshals the index file.
func LoadIndex(path string) (*Index, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index: %w", err)
	}

	return &index, nil
}