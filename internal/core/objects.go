package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// Object is the interface for all Zark objects (blob, tree, commit).
type Object interface {
	Hash() string
	Type() string
	Data() []byte
}

// Blob represents the content of a file.
type Blob struct {
	Content []byte
	hash    string
}

func NewBlob(content []byte) *Blob {
	hash := sha256.Sum256(content)
	return &Blob{
		Content: content,
		hash:    hex.EncodeToString(hash[:]),
	}
}

func (b *Blob) Hash() string { return b.hash }
func (b *Blob) Type() string { return "blob" }
func (b *Blob) Data() []byte { return b.Content }

// Tree represents a directory structure.
type Tree struct {
	Entries []TreeEntry
	hash    string
}

type TreeEntry struct {
	Mode string `json:"mode"`
	Name string `json:"name"`
	Hash string `json:"hash"`
	Type string `json:"type"`
}

func NewTree(entries []TreeEntry) *Tree {
	data, _ := json.Marshal(entries)
	hash := sha256.Sum256(data)
	return &Tree{
		Entries: entries,
		hash:    hex.EncodeToString(hash[:]),
	}
}

func (t *Tree) Hash() string { return t.hash }
func (t *Tree) Type() string { return "tree" }
func (t *Tree) Data() []byte {
	data, _ := json.Marshal(t.Entries)
	return data
}

// Commit represents a snapshot of the repository at a specific time.
type Commit struct {
	TreeHash  string    `json:"tree"`
	Parent    string    `json:"parent,omitempty"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	hash      string
}

func NewCommit(treeHash, parent, author, email, message string) *Commit {
	commit := &Commit{
		TreeHash:  treeHash,
		Parent:    parent,
		Author:    author,
		Email:     email,
		Timestamp: time.Now(),
		Message:   message,
	}

	data, _ := json.Marshal(commit)
	hash := sha256.Sum256(data)
	commit.hash = hex.EncodeToString(hash[:])

	return commit
}

func (c *Commit) Hash() string { return c.hash }
func (c *Commit) Type() string { return "commit" }
func (c *Commit) Data() []byte {
	data, _ := json.Marshal(c)
	return data
}