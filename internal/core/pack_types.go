package core

// PackEntry represents an object's metadata within the pack index.
type PackEntry struct {
	Hash   string
	Offset uint64
	CRC    uint32
}

// Constants for object types in packfiles, mirroring Git's format.
const (
	_             = iota // 0 is not used
	OBJ_COMMIT    = 1    // A commit object
	OBJ_TREE      = 2    // A tree object (directory listing)
	OBJ_BLOB      = 3    // A blob object (file content)
	_             = 4    // Reserved
	_             = 5    // Reserved
	OBJ_OFS_DELTA = 6    // A delta against an object at a specific offset in the pack
	OBJ_REF_DELTA = 7    // A delta against an object identified by its hash
)