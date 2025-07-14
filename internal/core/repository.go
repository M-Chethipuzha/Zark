package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Repository represents a Zark repository's structure and paths.
type Repository struct {
	Path       string
	ZarkDir    string
	ObjectsDir string
	RefsDir    string
	ConfigPath string
	HeadPath   string
	IndexPath  string
}

// Config holds the repository's configuration.
type Config struct {
	User UserConfig `json:"user"`
	Core CoreConfig `json:"core"`
}

type UserConfig struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CoreConfig struct {
	Bare bool `json:"bare"`
}

// NewRepository creates a new Repository struct for a given path.
func NewRepository(path string) *Repository {
	zarkDir := filepath.Join(path, ".zark")
	return &Repository{
		Path:       path,
		ZarkDir:    zarkDir,
		ObjectsDir: filepath.Join(zarkDir, "objects"),
		RefsDir:    filepath.Join(zarkDir, "refs"),
		ConfigPath: filepath.Join(zarkDir, "config"),
		HeadPath:   filepath.Join(zarkDir, "HEAD"),
		IndexPath:  filepath.Join(zarkDir, "index"),
	}
}

// Init initializes the directory structure and default files for a new repository.
func (r *Repository) Init() error {
	dirs := []string{
		r.ZarkDir,
		r.ObjectsDir,
		r.RefsDir,
		filepath.Join(r.RefsDir, "heads"),
		filepath.Join(r.RefsDir, "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create initial config
	config := Config{
		User: UserConfig{
			Name:  "User",
			Email: "user@example.com",
		},
		Core: CoreConfig{
			Bare: false,
		},
	}
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(r.ConfigPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Create initial HEAD pointing to the main branch
	if err := os.WriteFile(r.HeadPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return fmt.Errorf("failed to write HEAD: %w", err)
	}

	// Create empty index
	index := NewIndex()
	if err := index.Save(r.IndexPath); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

// Exists checks if a .zark directory exists at the repository path.
func (r *Repository) Exists() bool {
	_, err := os.Stat(r.ZarkDir)
	return err == nil
}

// GetConfig reads and unmarshals the repository's config file.
func (r *Repository) GetConfig() (*Config, error) {
	data, err := os.ReadFile(r.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}