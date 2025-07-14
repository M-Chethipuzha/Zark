package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// Storage handles reading from and writing to the object database.
type Storage struct {
	objectsDir string
}

func NewStorage(objectsDir string) *Storage {
	return &Storage{
		objectsDir: objectsDir,
	}
}

// Store writes an object to the database.
func (s *Storage) Store(obj Object) error {
	hash := obj.Hash()
	dir := filepath.Join(s.objectsDir, hash[:2])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create object directory: %w", err)
	}

	filename := filepath.Join(dir, hash[2:])
	if err := os.WriteFile(filename, obj.Data(), 0644); err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}

	return nil
}

// Load reads an object from the database.
func (s *Storage) Load(hash string) ([]byte, error) {
	filename := filepath.Join(s.objectsDir, hash[:2], hash[2:])
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read object %s: %w", hash, err)
	}
	return data, nil
}