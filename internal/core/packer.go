package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Packer handles the creation of packfiles and their indexes.
type Packer struct {
	repo         *Repository
	storage      *Storage
	looseObjects []string
	dmp          *diffmatchpatch.DiffMatchPatch
}

func NewPacker(repo *Repository) (*Packer, error) {
	var looseObjects []string
	err := filepath.Walk(repo.ObjectsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && len(info.Name()) == 62 {
			dir := filepath.Base(filepath.Dir(path))
			hash := dir + info.Name()
			looseObjects = append(looseObjects, hash)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan for loose objects: %w", err)
	}

	return &Packer{
		repo:         repo,
		storage:      NewStorage(repo),
		looseObjects: looseObjects,
		dmp:          diffmatchpatch.New(),
	}, nil
}

// PackObjects creates a single .pack file and a .idx file with delta compression.
func (p *Packer) PackObjects() (string, error) {
	if len(p.looseObjects) == 0 {
		return "", fmt.Errorf("no loose objects to pack")
	}

	sort.Strings(p.looseObjects)

	var packData bytes.Buffer
	var packEntries []PackEntry
	packedObjects := make(map[string][]byte)

	packData.Write([]byte("PACK"))
	binary.Write(&packData, binary.BigEndian, uint32(2))
	binary.Write(&packData, binary.BigEndian, uint32(len(p.looseObjects)))

	for _, hash := range p.looseObjects {
		objData, err := p.storage.loadLoose(filepath.Join(p.repo.ObjectsDir, hash[:2], hash[2:]), hash)
		if err != nil {
			return "", fmt.Errorf("could not load object %s for packing: %w", hash, err)
		}
		packedObjects[hash] = objData

		baseHash, delta := p.findBestDelta(hash, objData, packedObjects)

		currentOffset := uint64(packData.Len())
		packEntries = append(packEntries, PackEntry{Hash: hash, Offset: currentOffset})

		var objType uint8
		var dataToWrite []byte

		if delta != nil {
			objType = OBJ_REF_DELTA
			baseHashBytes, _ := hex.DecodeString(baseHash)
			dataToWrite = append(baseHashBytes, delta...)
		} else {
			objType = OBJ_BLOB
			dataToWrite = objData
		}

		p.writePackObjectHeader(&packData, objType, uint64(len(dataToWrite)))

		writer := zlib.NewWriter(&packData)
		writer.Write(dataToWrite)
		writer.Close()
	}

	packHashBytes := sha256.Sum256(packData.Bytes())
	packHash := hex.EncodeToString(packHashBytes[:])
	packDir := filepath.Join(p.repo.ZarkDir, "pack")
	os.MkdirAll(packDir, 0755)
	packFilePath := filepath.Join(packDir, fmt.Sprintf("pack-%s.pack", packHash))

	checksum := sha256.Sum256(packData.Bytes())
	packData.Write(checksum[:])

	if err := os.WriteFile(packFilePath, packData.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to write packfile: %w", err)
	}

	if err := p.writeIndex(packHash, packEntries); err != nil {
		return "", fmt.Errorf("failed to write pack index: %w", err)
	}

	for _, hash := range p.looseObjects {
		os.Remove(filepath.Join(p.repo.ObjectsDir, hash[:2], hash[2:]))
	}

	return packHash, nil
}

func (p *Packer) findBestDelta(currentHash string, currentData []byte, packedObjects map[string][]byte) (string, []byte) {
	var bestDelta []byte
	var bestBaseHash string

	for baseHash, baseData := range packedObjects {
		if baseHash == currentHash {
			continue
		}

		if len(baseData) < len(currentData)/2 || len(baseData) > len(currentData)*2 {
			continue
		}

		patches := p.dmp.PatchMake(string(baseData), string(currentData))
		deltaText := p.dmp.PatchToText(patches)

		if bestDelta == nil || len(deltaText) < len(bestDelta) {
			bestDelta = []byte(deltaText)
			bestBaseHash = baseHash
		}
	}

	if bestDelta != nil && len(bestDelta) < len(currentData) {
		return bestBaseHash, bestDelta
	}

	return "", nil
}

func (p *Packer) writePackObjectHeader(buf *bytes.Buffer, objType uint8, size uint64) {
	headerByte := (objType << 4) | uint8(size&0x0F)
	size >>= 4
	var header []byte
	for {
		if size == 0 {
			header = append(header, headerByte)
			break
		}
		headerByte |= 0x80
		header = append(header, headerByte)
		headerByte = uint8(size & 0x7F)
		size >>= 7
	}
	for i, j := 0, len(header)-1; i < j; i, j = i+1, j-1 {
		header[i], header[j] = header[j], header[i]
	}
	buf.Write(header)
}

func (p *Packer) GetObjectCount() int {
	return len(p.looseObjects)
}

func (p *Packer) writeIndex(packHash string, entries []PackEntry) error {
	var indexData bytes.Buffer

	indexData.Write([]byte{0xff, 't', 'O', 'c'})
	binary.Write(&indexData, binary.BigEndian, uint32(2))

	var fanout [256]uint32
	for i := range fanout {
		count := 0
		for _, entry := range entries {
			hashByte, _ := hex.DecodeString(entry.Hash[:2])
			if hashByte[0] <= byte(i) {
				count++
			}
		}
		fanout[i] = uint32(count)
	}
	binary.Write(&indexData, binary.BigEndian, fanout)

	for _, entry := range entries {
		hashBytes, _ := hex.DecodeString(entry.Hash)
		indexData.Write(hashBytes)
	}

	for range entries {
		binary.Write(&indexData, binary.BigEndian, uint32(0))
	}

	for _, entry := range entries {
		binary.Write(&indexData, binary.BigEndian, uint32(entry.Offset))
	}

	packChecksum, _ := hex.DecodeString(packHash)
	indexData.Write(packChecksum)

	indexChecksum := sha256.Sum256(indexData.Bytes())
	indexData.Write(indexChecksum[:])

	indexFilePath := filepath.Join(p.repo.ZarkDir, "pack", fmt.Sprintf("pack-%s.idx", packHash))
	return os.WriteFile(indexFilePath, indexData.Bytes(), 0644)
}