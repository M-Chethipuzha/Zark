package core

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Storage handles reading from and writing to the object database.
type Storage struct {
	repo *Repository
	dmp  *diffmatchpatch.DiffMatchPatch
}

func NewStorage(repo *Repository) *Storage {
	return &Storage{
		repo: repo,
		dmp:  diffmatchpatch.New(),
	}
}

// Store compresses and writes an object to the database as a loose object.
func (s *Storage) Store(obj Object) error {
	hash := obj.Hash()
	dir := filepath.Join(s.repo.ObjectsDir, hash[:2])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create object directory: %w", err)
	}

	var compressedData bytes.Buffer
	writer := zlib.NewWriter(&compressedData)
	_, err := writer.Write(obj.Data())
	if err != nil {
		return fmt.Errorf("failed to compress object data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close zlib writer: %w", err)
	}

	filename := filepath.Join(dir, hash[2:])
	return os.WriteFile(filename, compressedData.Bytes(), 0644)
}

// Load reads and decompresses an object from the database,
// checking loose objects first, then packfiles.
func (s *Storage) Load(hash string) ([]byte, error) {
	loosePath := filepath.Join(s.repo.ObjectsDir, hash[:2], hash[2:])
	if _, err := os.Stat(loosePath); err == nil {
		return s.loadLoose(loosePath, hash)
	}

	return s.loadFromPack(hash)
}

func (s *Storage) loadLoose(path, hash string) ([]byte, error) {
	compressedData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read object %s: %w", hash, err)
	}

	b := bytes.NewReader(compressedData)
	reader, err := zlib.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader for object %s: %w", hash, err)
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func (s *Storage) loadFromPack(hash string) ([]byte, error) {
	packDir := filepath.Join(s.repo.ZarkDir, "pack")
	var data []byte
	var found bool

	err := filepath.Walk(packDir, func(path string, info os.FileInfo, err error) error {
		if found || err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".idx") {
			packData, err := s.findInPack(path, hash)
			if err == nil {
				data = packData
				found = true
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("object not found in loose objects or packfiles: %s", hash)
	}
	return data, nil
}

func (s *Storage) findInPack(idxPath, targetHash string) ([]byte, error) {
	idxData, err := os.ReadFile(idxPath)
	if err != nil {
		return nil, err
	}

	shaOffset := 8 + 256*4
	numObjects := (len(idxData) - shaOffset - 20 - 20) / (32 + 4)

	for i := 0; i < numObjects; i++ {
		hashStart := shaOffset + i*32
		hashInIndex := string(idxData[hashStart : hashStart+32])

		if hashInIndex == targetHash {
			offsetStart := shaOffset + numObjects*32 + i*4
			offset := binary.BigEndian.Uint32(idxData[offsetStart : offsetStart+4])

			packPath := strings.TrimSuffix(idxPath, ".idx") + ".pack"
			packFile, err := os.Open(packPath)
			if err != nil {
				return nil, err
			}
			defer packFile.Close()

			packFile.Seek(int64(offset), 0)

			r, err := zlib.NewReader(packFile)
			if err != nil {
				return nil, err
			}
			defer r.Close()

			return io.ReadAll(r)
		}
	}

	return nil, fmt.Errorf("hash not found in index %s", idxPath)
}